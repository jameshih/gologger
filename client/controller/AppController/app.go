package AppController

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AppController struct {
	beego.Controller
}

func (p *AppController) Index() {

	logs.Debug("enter app controller")
	p.TplName = "index/index.html"
}
