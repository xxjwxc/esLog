package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/xxjwxc/esLog/view/es"
	"github.com/xxjwxc/esLog/view/oplogger"
)

type Logger struct {
}

//添加一个日志
func (s *Logger) AddLog() {
	es.Url = "http://192.168.198.17:9200/"
	es.Index = "wms_log"

	e := es.GetClient()

	var eslog es.ESLog
	eslog.Topic = "topic"
	eslog.EType = oplogger.EOpType_EOpGunbuster
	eslog.UserName = "username"
	eslog.Ekey = "iddd"
	eslog.ELevel = oplogger.ELogLevel_EOperate
	eslog.Desc = "desc yes haha oye"
	eslog.Attach = "attach"
	eslog.CreatTime = time.Now()

	var tmp []interface{}
	for i := 0; i < 3; i++ {
		eslog.Ekey = fmt.Sprintf("iddd-%v", i)
		eslog.Desc = fmt.Sprintf("%v-%v", eslog.Desc, i)
		tmp = append(tmp, eslog)
	}

	b := e.BulkAdd(es.Index, es.Index, "", tmp)
	if !b {
		fmt.Println(e.Err)
	}
}

//搜索日志
func (s *Logger) Search() {
	es.Url = "http://192.168.198.17:9200/"
	es.Index = "wms_log"

	//精确搜索
	term := make(map[string]interface{})
	term["topic"] = "topic"
	term["etype"] = oplogger.EOpType_EOpGunbuster
	term["user_name"] = "username"
	term["ekey"] = "iddd-1"
	term["elevel"] = oplogger.ELogLevel_EOperate
	//模糊匹配
	match := make(map[string]interface{})
	match["desc"] = "desc"
	match["attach"] = "attach"

	timeCase := make(map[string]es.CaseSection)
	timeCase["creat_time"] = es.CaseSection{
		Min: time.Now().AddDate(0, 0, -10),
		Max: time.Now(),
	}

	//构造搜索器
	var que es.EsQuery
	que.OnPages(0, 10).OnTerm(term).OnMatch(match).OnRangeTime(timeCase)
	data1, _ := json.Marshal(que.OnSource())
	fmt.Println(string(data1))

	client := es.GetClient()
	var eslog []es.ESLog
	client.Search(es.Index, es.Index,
		que.OnSource(), func(e []byte) error {
			var tmp es.ESLog
			err := json.Unmarshal(e, &tmp)
			if err != nil {
				log.Println(err)
			} else {
				eslog = append(eslog, tmp)
			}
			return err
		})

	fmt.Println(eslog)
}
