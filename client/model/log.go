package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jameshih/gologger/tailf"
	"go.etcd.io/etcd/clientv3"
)

type LogInfo struct {
	AppId      int    `db:"app_id"`
	AppName    string `db:"app_name"`
	LogId      int    `db:"log_id"`
	Statue     int    `db:"status"`
	CreateTime string `db:"create_time"`
	LogPath    string `db:"log_path"`
	Topic      string `db:"topic"`
}

var (
	etcdClient *clientv3.Client
)

func InitEtcd(cli *clientv3.Client) {
	etcdClient = cli
}

func GetAllLogInfo() (loglist []LogInfo, err error) {
	err = Db.Select(&loglist,
		"select a.app_id, b.app_name, a.log_id, a.status, a.create_time, a.log_path, a.topic from tbl_log_info a, tbl_app_info b where a.app_id=b.app_id")
	if err != nil {
		logs.Warn("exec failed, ", err)
		return
	}
	return
}

func CreateLog(info *LogInfo) (err error) {
	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("create log failed Db.Begin error: %v", err)
		return
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}
		conn.Commit()
	}()
	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", info.AppName)
	if err != nil || len(appId) == 0 {
		logs.Warn("get app_id failed, error: %v", err)
		return
	}

	info.AppId = appId[0]
	r, err := conn.Exec("insert into tbl_log_info(app_id, log_path, topic)values(?,?,?)", info.AppId, info.LogPath, info.Topic)
	if err != nil {
		logs.Warn("create log failed Db.Exec error: %v", err)
		return
	}
	_, err = r.LastInsertId()
	if err != nil {
		logs.Warn("create log failed, Db.LastInsertId error: %v", err)
		return
	}
	return
}

func SetLogConfigToEtcd(etcdKey string, info *LogInfo) (err error) {
	var logConfArr []tailf.CollectConf
	logConfArr = append(logConfArr, tailf.CollectConf{
		LogPath: info.LogPath,
		Topic:   info.Topic,
	})

	data, err := json.Marshal(logConfArr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	_, err = etcdClient.Put(ctx, etcdKey, string(data))
	cancel()
	if err != nil {
		switch err {
		case context.Canceled:
			logs.Warn("ctx is canceled by another routine: %v", err)
		case context.DeadlineExceeded:
			logs.Warn("ctx is attached with a deadline is exceeded: %v", err)
		case rpctypes.ErrEmptyKey:
			logs.Warn("client-side error: %v", err)
		default:
			logs.Warn("bad cluster endpoints, which are not etcd servers: %v", err)
		}
	}
	logs.Debug("put etcd succ, data: %v", string(data))
	return
}
