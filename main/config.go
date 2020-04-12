package main

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/config"
	"github.com/jameshih/gologger/tailf"
)

var (
	appConfig *Config
)

type Config struct {
	logLevel string
	logPath  string

	chanSize    int
	kafkaAddr   string
	kafkaPort   int
	collectConf []tailf.CollectConf
}

func loadCollectConf(conf config.Configer) (err error) {
	var cc tailf.CollectConf
	cc.LogPath = conf.String("collect::log_path")
	if len(cc.LogPath) == 0 {
		err = errors.New("invalid collect::log_path")
		return
	}
	cc.Topic = conf.String("collect::topic")
	if len(cc.Topic) == 0 {
		err = errors.New("invalid collect::topic")
		return
	}

	appConfig.collectConf = append(appConfig.collectConf, cc)
	return
}

func initConfig(configType, filename string) (err error) {
	conf, err := config.NewConfig(configType, filename)
	if err != nil {
		fmt.Println("new configg failed, err:", err)
		return
	}
	appConfig = &Config{}
	appConfig.logLevel = conf.String("logs::log_level")

	// set default val for logLevel is not define in config file
	if len(appConfig.logLevel) == 0 {
		appConfig.logLevel = "debug"
	}

	appConfig.logPath = conf.String("logs::log_path")

	// set default val for logPath is not define in config file
	if len(appConfig.logPath) == 0 {
		appConfig.logPath = "./logs"
	}

	appConfig.chanSize, err = conf.Int("collect::chan_size")
	if err != nil {
		appConfig.chanSize = 100
	}

	appConfig.kafkaAddr = conf.String("kafka::server_ip")
	if len(appConfig.kafkaAddr) == 0 {
		err = fmt.Errorf("invalid kafka address")
		return
	}

	appConfig.kafkaPort, err = conf.Int("kafka::server_port")
	if err != nil {
		err = fmt.Errorf("invalid kafka port, err:%v", err)
		return
	}

	err = loadCollectConf(conf)
	if err != nil {
		fmt.Printf("load collect conf failed, err:%v\n", err)
		return
	}
	return
}
