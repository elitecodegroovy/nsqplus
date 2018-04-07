package nsq

import (
	"encoding/json"
	"fmt"
	"time"
)



type ConsumerWorkerHandler struct {
	q                *Consumer
	messagesReceived int
	messagesFailed   int
}


func (h *ConsumerWorkerHandler) LogFailedMessage(message *Message) {
	h.messagesFailed++
	fmt.Sprintf("message failed [%d]", h.messagesFailed)
}

func (h *ConsumerWorkerHandler) HandleMessage(message *Message) error {
	data := struct {
		Msg string
	}{}

	err := json.Unmarshal(message.Body, &data)
	if err != nil {
		return err
	}

	msg := data.Msg
	h.messagesReceived++
	fmt.Println("message:", string(msg), ", message count", h.messagesReceived)

	return nil
}

func ReadNsqMessage(topicName string, cb func(c *Config)){
	config := NewConfig()
	// so that the test can simulate reaching max requeues and a call to LogFailedMessage
	config.DefaultRequeueDelay = 0
	// so that the test wont timeout from backing off
	config.MaxBackoffDuration = time.Millisecond * 50
	if cb != nil {
		cb(config)
	}
	topicName = topicName + time.Now().Local().Format("2006-01-02")
	q, _ := NewConsumer(topicName, "ch", config)
	// q.SetLogger(nullLogger, LogLevelInfo)

	h := &ConsumerWorkerHandler{
		q: q,
	}
	q.AddHandler(h)

	err := q.ConnectToNSQD(addr)
	if err != nil {
		fmt.Errorf("can't connect NSQD process ,error %s", err)
		<-q.StopChan
	}
}