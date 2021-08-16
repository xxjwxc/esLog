package logerlogic

import "github.com/xxjwxc/esLog/view/oplogger"

////////////////////////////////////////////////////////////////////////////
///
////////////////////////////////////////////////////////////////////////////

//三元组日志 操作人、操作时间、操作内容
type OpLogerTuple struct {
	LogerBase

	//ILogerTrimAdd
}

//添加一个
/*
e:业务类型
userName用户信息
*/
func (o *OpLogerTuple) AddOne(e int32, topic, userName, key, desc, attach string, l oplogger.ELogLevel) {
	var info LogerInfo
	info.CreatLoger(e, topic)
	info.SetUserName(userName).SetKey(key).SetDesc(desc)
	info.SetLevel(l)
	o.Add(&info.Loger)
}

//添加接口
func (o *OpLogerTuple) GetListParms() []*oplogger.LogerInfo {
	return o.Loger
}
