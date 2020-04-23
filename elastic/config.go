package main

import (
	"fmt"

	"github.com/astaxie/beego/config"
)

type LogConfig struct {
	KafkaAddr  string
	KafkaPort  int
	KafkaTopic string
	ESAddr     string
	ESPort     int
	LogPath    string
	LogLevel   string
}

var (
	logConfig *LogConfig
)

func initConfig(confType string, filename string) (err error) {
	conf, err := config.NewConfig(confType, filename)
	if err != nil {
		fmt.Println("new config failed, err:", err)
		return
	}
	logConfig = &LogConfig{}
	logConfig.LogLevel = conf.String("logs::log_level")

	// set default val for LogLevel is not define in config file
	if len(logConfig.LogLevel) == 0 {
		logConfig.LogLevel = "debug"
	}

	logConfig.LogPath = conf.String("logs::log_path")

	// set default val for LogPath is not define in config file
	if len(logConfig.LogPath) == 0 {
		logConfig.LogPath = "./logs"
	}

	logConfig.KafkaAddr = conf.String("kafka::server_ip")
	if len(logConfig.KafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka address")
		return
	}

	logConfig.KafkaPort, err = conf.Int("kafka::server_port")
	if err != nil {
		err = fmt.Errorf("invalid kafka port, err:%v", err)
		return
	}

	logConfig.KafkaTopic = conf.String("kafka::topic")
	if len(logConfig.KafkaTopic) == 0 {
		err = fmt.Errorf("invalid topic, err:%v", err)
		return
	}

	logConfig.ESAddr = conf.String("es::es_ip")
	if len(logConfig.ESAddr) == 0 {
		err = fmt.Errorf("invalid es ip, err:%v", err)
		return
	}

	logConfig.ESPort, err = conf.Int("es::es_port")
	if err != nil {
		err = fmt.Errorf("invalid es port, err:%v", err)
		return
	}

	return
}
