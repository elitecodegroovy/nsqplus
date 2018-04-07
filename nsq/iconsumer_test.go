package nsq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
	"crypto/tls"
	"nsqplus"
)

type IHandler struct {
	t                *testing.T
	q                *Consumer
	messagesSent     int
	messagesReceived int
	messagesFailed   int
}

func TestIConsumerSnappy(t *testing.T) {
	consumerITest(t, func(c *Config) {
		c.Snappy = true
	})
}

func TestIConsumerTLSDeflate(t *testing.T) {
	consumerITest(t, func(c *Config) {
		c.TlsV1 = true
		c.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		c.Deflate = true
	})
}

func (h *IHandler) LogFailedMessage(message *Message) {
	h.messagesFailed++
	h.q.Stop()
}

func (h *IHandler) HandleMessage(message *Message) error {
	data := struct {
		Msg string
	}{}

	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		return err
	}

	msg := data.Msg
	fmt.Println("message:", string(msg))
	h.messagesReceived++
	return nil
}



func TestIConsumerTLSSnappy(t *testing.T) {
	consumerTest(t, func(c *Config) {
		c.TlsV1 = true
		c.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		c.Snappy = true
	})
}


func TestIConsumerTLSClientCert(t *testing.T) {
	envDl := os.Getenv("NSQ_DOWNLOAD")
	if strings.HasPrefix(envDl, "nsq-0.2.24") || strings.HasPrefix(envDl, "nsq-0.2.27") {
		t.Log("skipping due to older nsqd")
		return
	}
	fmt.Println("getEnv", envDl)
	cert, _ := tls.LoadX509KeyPair("./test/client.pem", "./test/client.key")
	consumerITest(t, func(c *Config) {
		c.TlsV1 = true
		c.TlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
	})
}

func TestIConsumerTLSClientCertViaSet(t *testing.T) {
	envDl := os.Getenv("NSQ_DOWNLOAD")
	if strings.HasPrefix(envDl, "nsq-0.2.24") || strings.HasPrefix(envDl, "nsq-0.2.27") {
		t.Log("skipping due to older nsqd")
		return
	}
	consumerITest(t, func(c *Config) {
		c.Set("tls_v1", true)
		c.Set("tls_cert", "./test/client.pem")
		c.Set("tls_key", "./test/client.key")
		c.Set("tls_insecure_skip_verify", true)
	})
}

func consumerITest(t *testing.T, cb func(c *Config)) {
	config := NewConfig()
	// so that the test can simulate reaching max requeues and a call to LogFailedMessage
	config.DefaultRequeueDelay = 0
	// so that the test wont timeout from backing off
	config.MaxBackoffDuration = time.Millisecond * 50
	if cb != nil {
		cb(config)
	}
	topicName := "test_receiving"
	if config.Deflate {
		topicName = topicName + "_deflate"
	} else if config.Snappy {
		topicName = topicName + "_snappy"
	}
	if config.TlsV1 {
		topicName = topicName + "_tls"
	}
	topicName = topicName + time.Now().Local().Format("2006-01-02")
	q, _ := NewConsumer(topicName, "ch", config)
	// q.SetLogger(nullLogger, LogLevelInfo)

	h := &IHandler{
		t: t,
		q: q,
	}
	q.AddHandler(h)

	sendMessage(t, 4151, topicName, "put", []byte(`{"msg":"single"}`))
	sendMessage(t, 4151, topicName, "mput", []byte("{\"msg\":\"double\"}\n{\"msg\":\"double\"}"))
	sendMessage(t, 4151, topicName, "put", []byte("TOBEFAILED"))
	h.messagesSent = 4

	addr := "10.50.115.16:4150"
	err := q.ConnectToNSQD(addr)
	if err != nil {
		t.Fatal(err)
	}

	stats := q.Stats()
	if stats.Connections == 0 {
		t.Fatal("stats report 0 connections (should be > 0)")
	}

	err = q.ConnectToNSQD(addr)
	if err == nil {
		t.Fatal("should not be able to connect to the same NSQ twice")
	}

	<-q.StopChan

	stats = q.Stats()
	if stats.Connections != 0 {
		t.Fatalf("stats report %d active connections (should be 0)", stats.Connections)
	}

	stats = q.Stats()
	fmt.Sprintf("message received %d, finished %d", stats.MessagesReceived, stats.MessagesFinished)

	fmt.Sprintf("handler messagesReceived %d, sent %d by os ,failed %d", h.messagesReceived, h.messagesSent, h.messagesFailed)
}

func sendMessage(t *testing.T, port int, topic string, method string, body []byte) {
	httpclient := &http.Client{}
	endpoint := fmt.Sprintf("http://10.50.115.16:%d/%s?topic=%s", port, method, topic)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	resp, err := httpclient.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
		return
	}
	fmt.Println("resp: ", resp)
	resp.Body.Close()
}
