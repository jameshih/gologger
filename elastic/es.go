package main

import (
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch"
)

type LogMessage struct {
	App     string
	Topic   string
	Message string
}

var (
	esClient *elasticsearch.Client
)

func initElastic(addr string, port int) (err error) {
	//cfg := elasticsearch.Config{
	//Addresses: []string{
	//fmt.Sprintf("http://%s:%d", addr, port),
	//},
	//Transport: &http.Transport{
	//MaxIdleConnsPerHost:   10,
	//ResponseHeaderTimeout: time.Millisecond,
	//DialContext:           (&net.Dialer{Timeout: time.Nanosecond}).DialContext,
	//TLSClientConfig: &tls.Config{
	//MinVersion: tls.VersionTLS11,
	//},
	//},
	//}
	//esClient, err = elasticsearch.NewClient(cfg)
	esClient, err = elasticsearch.NewDefaultClient()
	if err != nil {
		return
	}
	return
}

func sendToES(topic string, data []byte) (err error) {
	msg := LogMessage{
		Topic:   topic,
		Message: string(data),
	}

	d, err := json.Marshal(msg)
	if err != nil {
		return
	}

	// _, err = esapi.IndexRequest{
	// 	Index:  topic,
	// 	OpType: topic,
	// 	// DocumentID: fmt.Sprintf("%d", i),
	// 	Body:    strings.NewReader(string(d)),
	// 	Refresh: "true",
	// }.Do(context.Background(), esClient)
	_, err = esClient.Index(
		topic,
		strings.NewReader(string(d)),
		esClient.Index.WithRefresh("true"),
	)

	if err != nil {
		panic(err)
	}
	return
}
