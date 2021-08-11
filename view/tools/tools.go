package tools

import (
	"encoding/json"
	"log"
	"time"

	"github.com/xxjwxc/esLog/view/es"
	"github.com/xxjwxc/esLog/view/oplogger"
)

//req请求类型转es日志
func ConvertRe2ESLog(req []*oplogger.LogerInfo) []es.ESLog {
	resp := make([]es.ESLog, 0, len(req))
	for _, v := range req {
		resp = append(resp, es.ESLog{
			Topic:     v.Topic,
			EType:     v.EType,
			UserName:  v.UserName,
			Ekey:      v.Ekey,
			ELevel:    v.ELevel,
			Desc:      v.Desc,
			Attach:    v.Attach,
			CreatTime: time.Unix(v.CreatTime, 0),
		})
	}

	return resp
}

//req请求类型转es日志
func ConvertESLogS2Re(req []es.ESLog) []*oplogger.LogerInfo {
	resp := make([]*oplogger.LogerInfo, 0, len(req))
	for _, v := range req {
		resp = append(resp, &oplogger.LogerInfo{
			Topic:     v.Topic,
			EType:     v.EType,
			UserName:  v.UserName,
			Ekey:      v.Ekey,
			ELevel:    v.ELevel,
			Desc:      v.Desc,
			Attach:    v.Attach,
			CreatTime: v.CreatTime.Unix(),
		})
	}

	return resp
}

//搜索并返回
func Search(index_name, type_name string, term map[string]interface{}, match map[string]interface{}, timeCase map[string]es.CaseSection, page, limit int32) []es.ESLog {

	//构造搜索器
	var que es.EsQuery
	que.OnPages(page, limit).OnTerm(term).OnMatch(match).OnRangeTime(timeCase)
	// data1, _ := json.Marshal(que.OnSource())
	// fmt.Println(string(data1))

	client, _ := es.New()
	var eslog []es.ESLog
	client.WithOption(es.WithIndexName(index_name), es.WithTypeName(type_name)).Search(que.OnSource(), func(e []byte) error {
		var tmp es.ESLog
		err := json.Unmarshal(e, &tmp)
		if err != nil {
			log.Println(err)
		} else {
			eslog = append(eslog, tmp)
		}
		return err
	})

	return eslog
}
