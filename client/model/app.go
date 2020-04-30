package model

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/logs"
	"github.com/jmoiron/sqlx"
)

type AppInfo struct {
	AppId       int    `db:"app_id"`
	AppName     string `db:"app_name"`
	AppType     string `db:"app_type"`
	CreateTime  string `db:"create_time"`
	DevelopPath string `db:"develop_path"`
	IP          []string
}

var (
	Db *sqlx.DB
)

func init() {
	db, err := sqlx.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/app")
	if err != nil {
		logs.Error("failed to connect to mysql, err: %v", err)
		return
	}
	Db = db
}

func GetAllInfo() (appInfo []AppInfo, err error) {
	err = Db.Select(&appInfo, "select app_id, app_name, app_type, create_time from tbl_app_info")
	if err != nil {
		logs.Warn("exec failed, ", err)
		return
	}
	return
}

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
		_, err = conn.Exec("inset into tbl_app_ip(app_id, ip)values(?,?)", id, ip)
		if err != nil {
			logs.Warn("create app ip failed, conn.Exec ip  error: %v", err)
			return
		}
	}
	return
}
