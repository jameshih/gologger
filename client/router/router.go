package router

import (
	"github.com/astaxie/beego"
	"github.com/jameshih/gologger/client/controller/AppController"
)

func init() {
	beego.Router("/index", &AppController.AppController{}, "*:Index")
}
