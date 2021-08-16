package oplogger

import (
	"github.com/golang/protobuf/proto"
)

type ELogLevel int32

const (
	// 默认
	ELogLevel_EDefault ELogLevel = 0
	// 调试
	ELogLevel_EDebug ELogLevel = 1
	// 普通
	ELogLevel_ENomal ELogLevel = 2
	// 操作级别
	ELogLevel_EOperate ELogLevel = 3
	// 警告
	ELogLevel_EWarning ELogLevel = 4
	// 错误级别
	ELogLevel_EError ELogLevel = 5
	// 致命
	ELogLevel_EFatal ELogLevel = 6
)

var ELogLevel_name = map[int32]string{
	0: "EDefault",
	1: "EDebug",
	2: "ENomal",
	3: "EOperate",
	4: "EWarning",
	5: "EError",
	6: "EFatal",
}

var ELogLevel_value = map[string]int32{
	"EDefault": 0,
	"EDebug":   1,
	"ENomal":   2,
	"EOperate": 3,
	"EWarning": 4,
	"EError":   5,
	"EFatal":   6,
}

func (x ELogLevel) String() string {
	return proto.EnumName(ELogLevel_name, int32(x))
}

type LogerInfo struct {
	//应用/服务的标识： 用来确定日志产生的应用服务器的唯一标识(可以细分)
	Topic string `json:"topic"`
	//业务唯一标识
	EType int32 `json:"etype"`
	//用户信息
	UserName string `json:"user_name"`
	//关键值
	Ekey string `json:"ekey"`
	//事件等级
	ELevel ELogLevel `json:"elevel"`
	//备注
	Desc string `json:"desc"`

	//附加字段
	Attach string `json:"attach"`

	//创建时间
	CreatTime int64 `json:"creat_time"`
}
