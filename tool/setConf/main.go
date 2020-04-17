package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/jameshih/gologger/tailf"
	"go.etcd.io/etcd/clientv3"
)

const (
	EtcdKey = "/backend/logagent/config/10.0.1.6"
)

func SetLogConfigToEtcd() {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	defer cli.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("connect succ")
	var logConfArr []tailf.CollectConf
	logConfArr = append(logConfArr, tailf.CollectConf{
		LogPath: "./logs/logagent.log",
		Topic:   "logs",
	})

	logConfArr = append(logConfArr, tailf.CollectConf{
		LogPath: "/project/nginx/logs/error2.log",
		Topic:   "log_err",
	})

	data, err := json.Marshal(logConfArr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	// delete log testing
	cli.Delete(ctx, EtcdKey)
	fmt.Printf("deleted %s", EtcdKey)
	return

	_, err = cli.Put(ctx, EtcdKey, string(data))
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	fmt.Println("setting value to etcd...")

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, EtcdKey)
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			log.Fatalf("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			log.Fatalf("client-side error: %v", err)
		default:
			log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s: %s\n", ev.Key, ev.Value)
	}

}

func main() {
	SetLogConfigToEtcd()
}
