package main

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/kafka"
	"github.com/jameshih/gologger/tailf"
)

func main() {
	filename := "./conf/logagent.conf"
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

	collectConf, err := initEtcd(appConfig.etcdAddr, appConfig.etcdKey)
	if err != nil {
		logs.Error("init etcd failed , err:%v", err)
		return
	}
	logs.Debug("initialize etcd  succ")

	err = tailf.InitTail(collectConf, appConfig.chanSize)
	if err != nil {
		logs.Error("init tail failed, err:%v", err)
		return
	}

	err = kafka.InitKafka(appConfig.kafkaAddr, appConfig.kafkaPort)
	if err != nil {
		logs.Error("init kafka failed, err:%v", err)
		return
	}

	go func() {
		var counter int
		for {
			logs.Debug("testing count:%d", counter)
			counter++
			time.Sleep(time.Second * 1)
		}
	}()

	logs.Debug("initialize all succ")
	err = startServer()
	if err != nil {
		logs.Error("start server failed, err:%v", err)
		return
	}
	logs.Info("program exited")
}
