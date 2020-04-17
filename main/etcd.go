package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/jameshih/gologger/tailf"
	"go.etcd.io/etcd/clientv3"
)

type EtcdClient struct {
	client *clientv3.Client
	keys   []string
}

var (
	etcdClient *EtcdClient
	cfg        clientv3.Config
)

func initEtcd(addr string, key string) (collectConf []tailf.CollectConf, err error) {
	cfg = clientv3.Config{
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
	collectConf = getFromEtcd(cli, key)
	logs.Debug("log config is %v", collectConf)
	initEtcdWatcher()
	return
}

func getFromEtcd(cli *clientv3.Client, key string) (collectConf []tailf.CollectConf) {
	for _, ip := range localIPArray {
		etcdKey := fmt.Sprintf("%s%s", key, ip)
		etcdClient.keys = append(etcdClient.keys, etcdKey)
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
			}
		}
	}
	return
}

func initEtcdWatcher() {
	for _, key := range etcdClient.keys {
		go watchKey(key)
	}
}

func watchKey(key string) {
	cli, err := clientv3.New(cfg)
	defer cli.Close()
	if err != nil {
		logs.Error("initEtcd failed, err:", err)
		return
	}
	//logs.Debug("key:%s", key)
	for {
		rch := cli.Watch(context.Background(), key)

		var collectConf []tailf.CollectConf
		var getConfSucc = true
		fmt.Print(getConfSucc)
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s]'s config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err = json.Unmarshal(ev.Kv.Value, &collectConf)
					if err != nil {
						logs.Error("key[%s], unmarshal[%s], err:%v", err)
						getConfSucc = false
						continue
					}
				}
				logs.Debug("get config form etcd %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", collectConf)
				tailf.UpdateConfig(collectConf)
			}
		}
	}
}
