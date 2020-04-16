package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/tailf"
	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
}

var (
	etcdClient *EtcdClient
)

func initEtcd(add string, key string) (err error) {
	cfg := clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	defer cli.Close()
	if err != nil {
		logs.Error("initEtcd failed, err:", err)
		return
	}

	etcdClient = &EtcdClient{
		client: cli,
	}
	if strings.HasSuffix(key, "/") == false {
		key = key + "/"
	}

	var collectConf []tailf.CollectConf
	for _, ip := range localIPArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := cli.Get(ctx, etcdKey)
		if err != nil {
			logs.Error("client get from etcd failed, err:%v", err)
			continue
		}
		logs.Debug("resp form etcd: %v", resp.Kvs)
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdKey {
				err = json.Unmarshal(v.Value, &collectConf)
				if err != nil {
					logs.Error("unmarshal failed, err:%v", err)
					continue
				}
				logs.Debug("log config is %v", collectConf)
			}
		}
	}
	return
}
