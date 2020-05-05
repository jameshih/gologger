package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/swithek/dotsqlx"
)

var (
	db *sqlx.DB
)

func init() {
	database, err := sqlx.Open("mysql", "root:toor@tcp(127.0.0.1:3306)/")
	if err != nil {
		logs.Error("failed to connect to mysql, err: %v", err)
		return
	}
	db = database
	return
}

func main() {
	dot, err := dotsql.LoadFromFile("./tables.sql")
	if err != nil {
		panic(err)
	}
	dotx := dotsqlx.Wrap(dot)
	_, err = dotx.Exec(db, "create-app-database")
	if err != nil {
		panic(err)
	}
	_, err = dotx.Exec(db, "use-app-database")
	if err != nil {
		panic(err)
	}

	_, err = dotx.Exec(db, "create-app-info-table")
	if err != nil {
		panic(err)
	}
	_, err = dotx.Exec(db, "create-app-ip-table")
	if err != nil {
		panic(err)
	}
	_, err = dotx.Exec(db, "create-log-info-table")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
