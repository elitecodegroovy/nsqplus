package main

import (
	"fmt"
	"github.com/elitecodegroovy/nsqplus/nsq"
)


func doConsumerTask(){
	fmt.Println("--start to process topic message...")
	nsq.ReadNsqMessage("007john", func(c *nsq.Config) {
		c.Snappy = true
	})
}

//nsq consumer task entrance
//func main() {
//	doConsumerTask()
//}
