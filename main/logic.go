package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/kafka"
	"github.com/jameshih/gologger/tailf"
)

func startServer() (err error) {
	var msg *tailf.ChanMsg
	fmt.Println("serving logger")
	for {
		msg = tailf.GetOneLine()
		err = sendToKafka(msg)
		if err != nil {
			logs.Error("send to Kafka failed, err:%v", err)
			time.Sleep(time.Second)
			continue
		}
	}
}

func sendToKafka(msg *tailf.ChanMsg) (err error) {
	kafka.SendToKafka(msg.Msg, msg.Topic)
	return
}
