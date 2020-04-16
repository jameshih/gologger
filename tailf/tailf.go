package tailf

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

type TailObj struct {
	tail *tail.Tail
	conf CollectConf
}

type ChanMsg struct {
	Msg   string
	Topic string
}

type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *ChanMsg
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *ChanMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {
	if len(conf) == 0 {
		err = fmt.Errorf("inva;id config for log collect, conf:%v", conf)
		return
	}

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *ChanMsg, 100),
	}

	for _, v := range conf {
		obj := &TailObj{
			conf: v,
		}
		tails, errTail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen: true,
			Follow: true,
			//Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll:      true,
		})

		if err != nil {
			err = errTail
			return
		}

		obj.tail = tails
		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

		go readFromTail(obj)
	}
	return
}

// read tail obj and add to channel
func readFromTail(tailObj *TailObj) {

	var line *tail.Line
	var ok bool

	for {
		line, ok = <-tailObj.tail.Lines
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
	}
}
