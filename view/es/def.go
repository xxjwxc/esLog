package es

import (
	"sync"
	"time"

	// "github.com/olivere/elastic"

	"github.com/xxjwxc/esLog/view/oplogger"
)

var esc *MyElastic
var once sync.Once

const retry = 3 //链接重试次数
//const PAGE_MAX_NUM = 10 //每页显示数据量

// ESLog ...
type ESLog struct {
	//应用/服务的标识： 用来确定日志产生的应用服务器的唯一标识(可以细分)
	Topic string `json:"topic"`
	//业务唯一标识
	EType int32 `json:"etype"`
	//用户信息
	UserName string `json:"user_name"`
	//关键值
	Ekey string `json:"ekey"`
	//事件等级
	ELevel oplogger.ELogLevel `json:"elevel"`
	//备注
	Desc string `json:"desc"`

	//附加字段
	Attach string `json:"attach"`

	//创建时间
	CreateTime time.Time `json:"create_time"`
}

// ESLogColumns 应用列名
var ESLogColumns = struct {
	Topic      string
	EType      string
	UserName   string
	Ekey       string
	ELevel     string
	Desc       string
	Attach     string
	CreateTime string
}{
	Topic:      "topic",
	EType:      "etype",
	UserName:   "user_name",
	Ekey:       "ekey",
	ELevel:     "elevel",
	Desc:       "desc",
	Attach:     "attach",
	CreateTime: "create_time",
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"topic":{
				"type":"keyword"
			},
			"user_name":{
				"type":"keyword"
			},
			"etype":{
				"type":"keyword"
			},
			"ekey":{
				"type":"keyword"
			},
			"elevel":{
				"type":"keyword"
			},
			"desc":{
				"type":"text",
				"fields": {
					"keyword": {
					  "type": "keyword",
					  "ignore_above": 256
					}
				},
				"analyzer": "ik_smart",		  
				"store": true,
				"fielddata": true
			},
			"attach":{
				"type":"text",
				"store": true,
				"fielddata": true
			},
			"create_time":{
				"type":"date"
			}
		}
	}
}`

//区间搜索使用
type CaseSection struct {
	Min interface{}
	Max interface{}
}

// type ESLog struct {

// 	//用来追踪一个请求的全服务调用流向
// 	TraceID string `json:"trace_id"`

// 	//应用/服务的唯一标识： 用来确定日志产生的应用服务器的唯一标识(可以细分)
// 	Topic  string `json:"topic"`
// 	Bundle string `json:"bundle"`
// 	Pid    string `json:"pid"`

// 	//用户信息
// 	UserID   string `json:"user_id"`
// 	UserName string `json:"user_name"`

// 	//业务唯一标识
// 	IdType oplogger.EOpType `json:"id_type"`
// 	Id     string           `json:"id"`
// 	SubID  string           `json:"sub_id"`

// 	//时间列表
// 	BeginTime time.Time `json:"begin_time"`
// 	EndTime   time.Time `json:"end_time"`
// 	CreatTime time.Time `json:"create_time"`

// 	//事件列表，描述
// 	ELevel oplogger.ELogLevel `json:"elevel"`
// 	//类型，1：自动，2：手动
// 	EType oplogger.EOpType `json:"etype"`
// 	EDesc string           `json:"edesc"`

// 	//变化值序列
// 	IType int32 `json:"itype"`
// 	//变化前值
// 	IOriginal string `json:"ioriginal"`
// 	//变化后值
// 	ITranslation string `json:"itranslation"`
// 	//变化值
// 	IVaiable string `json:"ivaiable"`

// 	//备注
// 	Desc string `json:"desc"`

// 	//预留字段
// 	DString string    `json:"dstring"`
// 	DInt    int64     `json:"dint"`
// 	DDate   time.Time `json:"ddate"`
// }

// const mapping = `
// {
// 	"settings":{
// 		"number_of_shards": 1,
// 		"number_of_replicas": 0
// 	},
// 	"mappings":{
// 		"oplogger":{
// 			"properties":{
// 				"trace_id":{
// 					"type":"keyword"
// 				},
// 				"topic":{
// 					"type":"keyword"
// 				},
// 				"bundle":{
// 					"type":"keyword"
// 				},
// 				"pid":{
// 					"type":"keyword"
// 				},
// 				"user_id":{
// 					"type":"keyword"
// 				},
// 				"user_name":{
// 					"type":"keyword"
// 				},
// 				"id_type":{
// 					"type":"keyword"
// 				},
// 				"id":{
// 					"type":"keyword"
// 				},
// 				"sub_id":{
// 					"type":"keyword"
// 				},
// 				"begin_time":{
// 					"type":"date"
// 				},
// 				"end_time":{
// 					"type":"date"
// 				},
// 				"create_time":{
// 					"type":"date"
// 				},
// 				"elevel":{
// 					"type":"keyword"
// 				},
// 				"etype":{
// 					"type":"keyword"
// 				},
// 				"edesc":{
// 					"type":"text",
// 					"store": true,
// 					"fielddata": true
// 				},
// 				"itype":{
// 					"type":"keyword"
// 				},
// 				"ioriginal":{
// 					"type":"keyword"
// 				},
// 				"itranslation":{
// 					"type":"keyword"
// 				},
// 				"ivaiable":{
// 					"type":"keyword"
// 				},
// 				"desc":{
// 					"type":"text",
// 					"store": true,
// 					"fielddata": true
// 				},
// 				"dstring":{
// 					"type":"text",
// 					"store": true,
// 					"fielddata": true
// 				},
// 				"dint":{
// 					"type":"keyword"
// 				},
// 				"ddate":{
// 					"type":"date"
// 				}
// 			}
// 		}
// 	}
// }`
