package main

import (
	"fmt"
	"go_dev/projects/log_collection/kafka"
	"go_dev/projects/log_collection/tailf"

	"github.com/astaxie/beego/logs"
)

func main() {
	filename := "./logagent.conf"
	err := initConfig("ini", filename)
	if err != nil {
		fmt.Printf("load conf failed, err:%v\n", err)
		panic("load conf failed")
	}

	err = initLogger()
	if err != nil {
		fmt.Printf("load loggger failed, err:%v\n", err)
		panic("load logger failed")
	}

	logs.Debug("load conf succ, config:%v", appConfig)

	err = tailf.InitTail(appConfig.collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}

	err = kafka.InitKafka(appConfig.kafkaAddr, appConfig.kafkaPort)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	logs.Debug("initialize all succ")
	err = startServer()
	if err != nil {
		logs.Error("start server failed, err:%v", err)
		return
	}
	logs.Info("program exited")
}
