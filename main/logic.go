package main

import (
	"go_dev/projects/log_collection/kafka"
	"go_dev/projects/log_collection/tailf"
	"time"

	"github.com/astaxie/beego/logs"
)

func startServer() (err error) {
	var msg *tailf.ChanMsg
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
