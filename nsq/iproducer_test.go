package nsq

import (
	"testing"
	"log"
)


func TestSimpleProducer1(t *testing.T){
	//init default config
	config := NewConfig()

	w, _ := NewProducer("10.50.115.16:4150", config)
	err := w.Ping()
	if err != nil {
		log.Fatalln("error ping 10.50.115.16:4150", err)
		// switch the second nsq. You can use nginx or HAProxy for HA.
		w, _ = NewProducer("10.50.115.17:4150", config)
	}
	w.SetLogger(nullLogger, LogLevelInfo)

	err2 := w.Publish("req-test", []byte("test"))
	if err2 != nil {
		panic("should lazily connect - %s")
	}

	w.Stop()
}
