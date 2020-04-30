package AppController

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/client/model"
)

type AppController struct {
	beego.Controller
}

func (p *AppController) AppList() {

	logs.Debug("enter app controller")
	p.Layout = "layout/layout.html"
	appList, err := model.GetAllInfo()
	if err != nil {
		// p.Data["Error"] = fmt.Sprintf("server busy...")
		p.Data["Error"] = err
		p.TplName = "app/error.html"
		logs.Warn("get app list failed, err: %v", err)
		return
	}
	logs.Debug("get app list succ, data:%v", appList)
	p.Data["appList"] = appList
	p.TplName = "app/index.html"
}

func (p *AppController) AppApply() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

func (p *AppController) AppCreate() {
	logs.Debug("enter index controller")
	appName := p.GetString("app_name")
	appType := p.GetString("app_type")
	developPath := p.GetString("develop_path")
	ipstr := p.GetString("iplist")

	if len(appName) == 0 || len(appType) == 0 || len(developPath) == 0 || len(ipstr) == 0 {
		p.Data["Error"] = fmt.Sprintf("invalid params")
		p.TplName = "app/error.html"
		logs.Warn("invalid params")
		return
	}

	appInfo := &model.AppInfo{}
	appInfo.AppName = appName
	appInfo.AppType = appType
	appInfo.DevelopPath = developPath
	appInfo.IP = strings.Split(ipstr, ",")

	_, err := model.InsertAppInfo(appInfo)
	if err != nil {
		p.Data["Error"] = fmt.Sprintf("failed to create new project, db busy")
		p.TplName = "app/error.html"
		logs.Warn("invalid params")
		return
	}

	p.Redirect("/app/list", 302)
}
