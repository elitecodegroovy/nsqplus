package nsqplus

import (
	"log"
	"os"
)

type Options struct {
	HTTPAddress string `flag:"http-address"`

	NSQDHTTPAddresses       []string `flag:"nsqd-http-address" cfg:"nsqd_http_addresses"`
	Logger logger
}


func NewOptions() *Options {
	return &Options{
		HTTPAddress:         "0.0.0.0:8800",
		NSQDHTTPAddresses:   []string{"10.50.115.16:4150"},
		Logger:              log.New(os.Stderr, "[nsqplus] ", log.Ldate|log.Ltime|log.Lmicroseconds),
	}
}

type ConsumerOptions struct {
	PushHttpAddress      string  `flag:"push-http-address"`

	NsqlookupHTTPAddress []string `flag:"nsqlookup-http-addresses" cfg:"nsqlookup_http_addresses"`
}

func NewConsumerOptions() *ConsumerOptions{
	return &ConsumerOptions{
		PushHttpAddress     : 	"http://10.50.115.14:8801/statistics/api/push",
		NsqlookupHTTPAddress:   []string{"10.50.115.16:4161"},
	}
}
