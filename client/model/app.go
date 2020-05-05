package model

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type AppInfo struct {
	AppId       int      `db:"app_id"`
	AppName     string   `db:"app_name"`
	AppType     string   `db:"app_type"`
	CreateTime  string   `db:"create_time"`
	DevelopPath string   `db:"develop_path"`
	IP          []string `db:"ip"`
}

var (
	Db *sqlx.DB
)

func InitDb(db *sqlx.DB) {
	Db = db
}

func GetAllInfo() (appInfo []AppInfo, err error) {
	err = Db.Select(&appInfo, "select app_id, app_name, app_type, create_time, develop_path from tbl_app_info")
	if err != nil {
		logs.Warn("get all app info failed, error: %v", err)
		return
	}
	return
}

func GetIPInfoByName(appName string) (iplist []string, err error) {
	var appId []int
	err = Db.Select(&appId, "select app_id from tbl_app_info where app_name=?", appName)
	if err != nil || len(appId) == 0 {
		logs.Warn("get app_id failed, error: %v", err)
		return
	}
	err = Db.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId[0])
	if err != nil {
		logs.Warn("get ip info by appid failed, error: %v", err)
		return
	}
	return
}

//func GetIPInfoById(appId int) (iplist []string, err error) {
//err = Db.Select(&iplist, "select ip from tbl_app_ip where app_id=?", appId)
//if err != nil {
//logs.Warn("get ip info by appid failed, error: %v", err)
//return
//}
//return
//}

func InsertAppInfo(appInfo *AppInfo) (id int64, err error) {
	conn, err := Db.Begin()
	if err != nil {
		logs.Warn("begin failed, err:%v", err)
		return
	}
	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}
		conn.Commit()
	}()
	r, err := conn.Exec("insert into tbl_app_info(app_name, app_type, develop_path)values(?, ?, ?)", appInfo.AppName, appInfo.AppType, appInfo.DevelopPath)
	if err != nil {
		logs.Warn("create app failed, Db.Exec error: %v", err)
		return
	}
	id, err = r.LastInsertId()
	if err != nil {
		logs.Warn("create app failed, Db.LastInsertId error: %v", err)
		return
	}
	for _, ip := range appInfo.IP {
		_, err = conn.Exec("insert into tbl_app_ip(app_id, ip)values(?,?)", id, ip)
		if err != nil {
			logs.Warn("create app ip failed, conn.Exec ip  error: %v", err)
			return
		}
	}
	return
}
