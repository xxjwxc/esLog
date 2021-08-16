package tools

import (
	"strconv"

	"github.com/xxjwxc/esLog/view/oplogger"
)

//参数判空
func IsEmpty(req []*oplogger.LogerInfo) bool {
	for _, v := range req {
		if CheckParam(v.Topic, v.UserName, v.Ekey, v.Desc, v.Attach) ||
			v.EType != 0 || v.ELevel != oplogger.ELogLevel_EDefault || v.CreatTime != 0 {
			return false
		}
	}

	return true
}

//检测参数
func CheckParam(params ...string) bool {
	for _, value := range params {
		if len(value) == 0 {
			return false
		}
	}
	return true
}

//
func Int64ToString(p int64) string {
	return strconv.FormatInt(p, 10)
}
