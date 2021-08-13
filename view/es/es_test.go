package es

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/xxjwxc/esLog/view/oplogger"
)

func TestEsAdd(m *testing.T) {
	url := "http://192.168.198.17:9200/"
	Index := "xxj_wms_log"

	es, _ := New(WithIndexName(Index), WithAddrs(url))

	var eslog ESLog
	eslog.Topic = "topic"
	eslog.EType = oplogger.EOpType_EOpGunbuster
	eslog.UserName = "username"
	eslog.Ekey = "iddd"
	eslog.ELevel = oplogger.ELogLevel_EOperate
	eslog.Desc = "desc yes haha oye"
	eslog.Attach = "attach"
	eslog.CreateTime = time.Now()

	var tmp []interface{}
	for i := 0; i < 3; i++ {
		eslog.Ekey = fmt.Sprintf("iddd-%v", i)
		eslog.Desc = fmt.Sprintf("%v-%v", eslog.Desc, i)
		tmp = append(tmp, eslog)
	}

	err := es.BulkAdd(tmp)
	if err != nil {
		fmt.Println(err)
	}

	// var eslogs []ESLog
	// es.Search(config.Config.ElasticSearch.Index,
	// 	config.Config.ElasticSearch.Index, `{"query":{"match_all":{}}}`, &eslogs)

	// fmt.Println(eslogs)
}

func TestSearch(t *testing.T) {
	url := "http://192.168.198.17:9200/"
	Index := "xxj_wms_log"

	es, _ := New(WithIndexName(Index), WithAddrs(url))

	// var ws []map[string]interface{}
	// b := es.SearchMap(config.Config.ElasticSearch.Index,
	// 	config.Config.ElasticSearch.Index, `{"query":{"match_all":{}}}`, &ws)
	// fmt.Println(b)
	// for _, v := range ws {
	// 	fmt.Println(v)
	// }
	// fmt.Println(len(ws))

	source := map[string]interface{}{
		"from": 0,
		"size": 10,
		"sort": map[string]interface{}{
			"creat_time": "desc",
		},
		"query": map[string]interface{}{
			"filtered": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"must": []interface{}{
							map[string]interface{}{
								"match": map[string]interface{}{
									"itype": 666,
								},
							},
							map[string]interface{}{
								"match": map[string]interface{}{
									"dint": 456,
								},
							},
						},
					},
				},
				"filter": map[string]interface{}{
					"and": []interface{}{
						map[string]interface{}{
							"range": map[string]interface{}{
								"creat_time": map[string]interface{}{
									"gte": time.Now().AddDate(0, 0, -1),
									"lte": time.Now(),
								},
							},
						},
						map[string]interface{}{
							"range": map[string]interface{}{
								"begin_time": map[string]interface{}{
									"gte": time.Now().AddDate(0, 0, -1),
									"lte": time.Now(),
								},
							},
						},
					},
				},
			},
		},
	}

	data1, _ := json.Marshal(source)
	fmt.Println(string(data1))
	//query := elastic.NewSearchSource().Query(elastic.NewMatchAllQuery()).From(0).Size(1)
	//`{"query":{"match_all":{}}}`
	var eslog []ESLog
	es.WithOption(WithIndexName(Index), WithTypeName(Index)).Search(source, func(e []byte) error {
		var tmp ESLog
		err := json.Unmarshal(e, &tmp)
		if err != nil {
			log.Println(err)
		} else {
			eslog = append(eslog, tmp)
		}
		return err
	})

	for _, v := range eslog {
		fmt.Println(v.CreateTime)
		fmt.Println("-----------------")
	}

}

func TestSearchObj(t *testing.T) {
	url := "http://192.168.198.17:9200/"
	Index := "xxj_wms_log"

	parm := make(map[string]interface{})
	parm["itype"] = 666
	parm["dint"] = 456

	timecase := make(map[string]CaseSection)
	timecase["creat_time"] = CaseSection{
		Min: time.Now().Unix(),
		Max: time.Now().AddDate(0, 0, -2).Unix(),
	}

	var que EsQuery
	que.OnPages(1, 1).OnMatch(parm).OnRangeTime(timecase)
	data1, _ := json.Marshal(que.OnSource())
	fmt.Println(string(data1))

	es, _ := New(WithIndexName(Index), WithAddrs(url))
	var eslog []ESLog
	es.WithOption(WithIndexName(Index), WithTypeName(Index)).Search(que.OnSource(), func(e []byte) error {
		var tmp ESLog
		err := json.Unmarshal(e, &tmp)
		if err != nil {
			log.Println(err)
		} else {
			eslog = append(eslog, tmp)
		}
		return err
	})

	for _, v := range eslog {
		fmt.Println(v.CreateTime)
		fmt.Println("-----------------")
	}

}

func TestTrackingOpLoger(t *testing.T) {
	url := "http://192.168.198.17:9200/"
	Index := "xxj_wms_log"

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

	timeCase := make(map[string]CaseSection)
	timeCase["creat_time"] = CaseSection{
		Min: time.Now().AddDate(0, 0, -1),
		Max: time.Now(),
	}

	//构造搜索器
	var que EsQuery
	que.OnPages(0, 10).OnTerm(term).OnMatch(match).OnRangeTime(timeCase)
	data1, _ := json.Marshal(que.OnSource())
	fmt.Println(string(data1))

	client, _ := New(WithIndexName(Index), WithAddrs(url))
	var eslog []ESLog
	client.WithOption(WithIndexName(Index), WithTypeName(Index)).Search(que.OnSource(), func(e []byte) error {
		var tmp ESLog
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
