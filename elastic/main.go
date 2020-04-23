package main

import (
	"github.com/astaxie/beego/logs"
)

func main() {
	err := initConfig("ini", "./conf/logagent.conf")
	if err != nil {
		panic(err)
	}

	err = initLogger(logConfig.LogPath, logConfig.LogLevel)
	if err != nil {
		panic(err)
	}
	logs.Debug("init logger succ")

	err = initKafka(logConfig.KafkaAddr, logConfig.KafkaPort, logConfig.KafkaTopic)
	if err != nil {
		logs.Error("initializing kafka failed, err: %v", err)
		return
	}
	logs.Debug("init kafka succ")

	err = initElastic(logConfig.ESAddr, logConfig.ESPort)
	if err != nil {
		logs.Error("initializing elasticsearch failed, err: %v", err)
		return
	}
	logs.Debug("init elasticsearch succ")

	err = startServer()
	if err != nil {
		logs.Error("start server failed, err: %v", err)
		return
	}
	logs.Warn("warning, program exited")
}
