package logerlogic

import "github.com/xxjwxc/esLog/view/oplogger"

// //日志 过滤器
// type ILogicFilter interface {
// 	//获取过滤之后的数据
// 	GetLogerFilter([]*oplogger.LogerInfo) interface{}
// 	//设置需要搜索的过滤器值
// 	SetLogerFilter(...interface{}) *oplogger.SearchReq
// }

//日志 添加 接口
type ILogerTrimAdd interface {
	//获取要添加的日志
	GetListParms() []*oplogger.LogerInfo
}
