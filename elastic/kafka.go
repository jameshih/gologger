package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type KafkaClient struct {
	client sarama.Consumer
	addr   string
	topic  string
	wg     sync.WaitGroup
}

var (
	kafkaClient *KafkaClient
)

func initKafka(addr string, port int, topic string) (err error) {
	kafkaClient = &KafkaClient{}
	consumer, err := sarama.NewConsumer(strings.Split(fmt.Sprintf("%s:%d", addr, port), ","), nil)
	if err != nil {
		return
	}

	kafkaClient.client = consumer
	kafkaClient.addr = addr
	kafkaClient.topic = topic
	return
}

func fetchFromKafka(pc sarama.PartitionConsumer, topic string, wg *sync.WaitGroup) {
	for msg := range pc.Messages() {
		logs.Debug("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		err := sendToES(topic, msg.Value)
		if err != nil {
			logs.Debug("send to es failed, err: %v", err)
		}
	}
	wg.Done()
}
