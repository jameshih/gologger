package main

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/client/model"
	_ "github.com/jameshih/gologger/client/router"
	"github.com/jmoiron/sqlx"
	"go.etcd.io/etcd/clientv3"
)

func initDb() (err error) {
	db, err := sqlx.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/app")
	if err != nil {
		logs.Error("failed to connect to mysql, err: %v", err)
		return
	}
	model.InitDb(db)
	return
}

func initEtcd() (err error) {
	cfg := clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		logs.Warn("initEtcd failed, error: %v", err)
		return
	}
	model.InitEtcd(cli)
	return
}

func main() {
	err := initDb()
	if err != nil {
		logs.Warn("initDb failed, error: %v", err)
		return
	}
	err = initEtcd()
	if err != nil {
		logs.Warn("initEtcd failed, error: %v", err)
		return
	}
	beego.Run()
}
