package main

import (
	"github.com/Shopify/sarama"
)

func startServer() (err error) {
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		return
	}

	for partition := range partitionList {
		pc, errRet := kafkaClient.client.ConsumePartition(kafkaClient.topic, int32(partition), sarama.OffsetNewest)
		if errRet != nil {
			err = errRet
			return
		}
		defer pc.AsyncClose()
		kafkaClient.wg.Add(1)
		fetchFromKafka(pc, kafkaClient.topic, &kafkaClient.wg)
	}
	kafkaClient.wg.Wait()
	return
}
