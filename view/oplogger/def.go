package oplogger

import "github.com/golang/protobuf/proto"

type EOpType int32

const (
	// 默认
	EOpType_EOpDefault EOpType = 0
	// 运单操作类日志
	EOpType_EOpTracking EOpType = 1
	// 采购操作类日志
	EOpType_EOpPurchase EOpType = 2
	// 采购下单操作类日志
	EOpType_EOpPurchaseOrder EOpType = 3
	// 商品作类日志
	EOpType_EOpProduct EOpType = 4
	// 腾云操作日志
	EOpType_EOpGunbuster EOpType = 5
	// 工单操作日志
	EOpType_EOpTicket EOpType = 6
	// 验货入库操作日志
	EOpType_EOpSTokin EOpType = 7
	// 打包封箱操作日志
	EOpType_EOpPicking EOpType = 8
)

var EOpType_name = map[int32]string{
	0: "EOpDefault",
	1: "EOpTracking",
	2: "EOpPurchase",
	3: "EOpPurchaseOrder",
	4: "EOpProduct",
	5: "EOpGunbuster",
	6: "EOpTicket",
	7: "EOpSTokin",
	8: "EOpPicking",
}
var EOpType_value = map[string]int32{
	"EOpDefault":       0,
	"EOpTracking":      1,
	"EOpPurchase":      2,
	"EOpPurchaseOrder": 3,
	"EOpProduct":       4,
	"EOpGunbuster":     5,
	"EOpTicket":        6,
	"EOpSTokin":        7,
	"EOpPicking":       8,
}

func (x EOpType) String() string {
	return proto.EnumName(EOpType_name, int32(x))
}

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
