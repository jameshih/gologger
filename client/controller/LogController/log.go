package LogController

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jameshih/gologger/client/model"
)

type LogController struct {
	beego.Controller
}

func (p *LogController) LogList() {

	logs.Debug("enter log controller")
	p.Layout = "layout/layout.html"
	loglist, err := model.GetAllLogInfo()
	if err != nil {
		// p.Data["Error"] = fmt.Sprintf("server busy...")
		p.Data["Error"] = err
		p.TplName = "layout/error.html"
		logs.Warn("get app list failed, err: %v", err)
		return
	}
	logs.Debug("get app list succ, data:%v", loglist)
	p.Data["loglist"] = loglist
	p.TplName = "log/index.html"
}

func (p *LogController) LogApply() {
	logs.Debug("enter index controller")
	p.Layout = "layout/layout.html"
	p.TplName = "log/apply.html"
}

func (p *LogController) LogCreate() {
	logs.Debug("enter index controller")
	appName := p.GetString("app_name")
	logPath := p.GetString("log_path")
	topic := p.GetString("topic")

	if len(appName) == 0 || len(logPath) == 0 || len(topic) == 0 {
		p.Data["Error"] = fmt.Sprintf("invalid params")
		p.TplName = "layout/error.html"
		logs.Warn("invalid params")
		return
	}

	logInfo := &model.LogInfo{}
	logInfo.AppName = appName
	logInfo.LogPath = logPath
	logInfo.Topic = topic

	err := model.CreateLog(logInfo)
	if err != nil {
		// p.Data["Error"] = fmt.Sprintf("failed to create new project, db busy")
		p.Data["Error"] = err
		p.TplName = "layout/error.html"
		logs.Warn("invalid params")
		return
	}

	p.Redirect("/log/list", 302)
}
