package logerlogic

import (
	"time"

	"github.com/xxjwxc/public/errors"

	"git.ezbuy.me/ezbuy/base/misc/context"
	"git.ezbuy.me/ezbuy/oplogger/rpc/oplogger"
)

type LogerInfo struct {
	Loger oplogger.LogerInfo
}

func (b *LogerInfo) CreatLoger(e oplogger.EOpType, topic string) *LogerInfo {
	//业务唯一标识
	b.Loger.EType = e
	b.Loger.Topic = topic

	b.Loger.CreatTime = time.Now().Unix()
	return b
}

//设置关键值
func (b *LogerInfo) SetKey(key string) *LogerInfo {
	b.Loger.Ekey = key
	return b
}

//设置用户信息
func (b *LogerInfo) SetUserName(userName string) *LogerInfo {
	b.Loger.UserName = userName
	return b
}

//事件等级
func (b *LogerInfo) SetLevel(l oplogger.ELogLevel) *LogerInfo {
	b.Loger.ELevel = l
	return b
}

//备注
func (b *LogerInfo) SetDesc(desc string) *LogerInfo {
	b.Loger.Desc = desc
	return b
}

//附加
func (b *LogerInfo) SetAttach(attach string) *LogerInfo {
	b.Loger.Attach = attach
	return b
}

////////////////////////////////////////////////////////////////////////////
///
////////////////////////////////////////////////////////////////////////////

type LogerBase struct {
	Loger []*oplogger.LogerInfo
}

func (l *LogerBase) Add(info *oplogger.LogerInfo) {
	l.Loger = append(l.Loger, info)
}

////////////////////////////////////////////////////////////////////////////
//////逻辑入口
////////////////////////////////////////////////////////////////////////////

// 添加日志到日志中心
func OnAddLogOplogger(ctx context.T, info ILogerTrimAdd) error {
	tmp := info.GetListParms()
	if len(tmp) == 0 {
		return ctx.Trace(errors.New("参数为空"))
	}

	go func() {
		oplogger.GetLogger().AddLog(ctx, &oplogger.AddLogReq{Info: tmp})
	}()

	return nil
}
