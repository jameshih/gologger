package router

import (
	"github.com/astaxie/beego"
	AppController "github.com/jameshih/gologger/client/controller/AppController"
)

func init() {
	beego.Router("/inde", &AppController.AppController{}, "*:index")
}
