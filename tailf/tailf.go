package tailf

import (
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int
	exitChan chan int
}

type ChanMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *ChanMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func InitTail(conf []CollectConf, chanSize int) (err error) {
	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *ChanMsg, 100),
	}

	if len(conf) == 0 {
		logs.Error("invalid config for log collect, conf:%v", conf)
		//err = fmt.Errorf("invalid config for log collect, conf:%v", conf)
		return
	}

	for _, v := range conf {
		createNewTask(v)
	}
	return
}

func GetOneLine() (msg *ChanMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

func UpdateConfig(confs []CollectConf) (err error) {
	// must have lock for multiple goroutines
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()
	for _, conf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs {
			if conf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}
		createNewTask(conf)
	}

	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, conf := range confs {
			if conf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}

		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}

	tailObjMgr.tailObjs = tailObjs
	return
}

func createNewTask(conf CollectConf) {
	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}
	tails, err := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})

	if err != nil {
		logs.Error("collect filename[%s] failed, err:%s", conf.LogPath, err)
		return
	}

	obj.tail = tails
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

	go readFromTail(obj)
}

// read tail obj and add to channel
func readFromTail(tailObj *TailObj) {

	var line *tail.Line
	var ok bool

	for {
		select {
		case line, ok = <-tailObj.tail.Lines:
			if !ok {
				logs.Warn("tail file close reopen, filename:%s\n", tailObj.tail.Filename)
				time.Sleep(100 * time.Microsecond)
				continue
			}
			chanMsg := &ChanMsg{
				Msg:   line.Text,
				Topic: tailObj.conf.Topic,
			}
			tailObjMgr.msgChan <- chanMsg
		case <-tailObj.exitChan:
			logs.Warn("tail obj will exit, conf:%v", tailObj.conf)
			return
		}
	}
}
