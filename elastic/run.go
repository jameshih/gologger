package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func startServer() (err error) {
	partitionList, err := kafkaClient.client.Partitions(kafkaClient.topic)
	if err != nil {
		return
	}

	fmt.Printf("Listening for kafka topic: %s on port: %d\n", logConfig.KafkaTopic, logConfig.KafkaPort)
	fmt.Printf("Publishing logs to ES on port: %d\nKibana serving on port 5601", logConfig.ESPort)
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
