package rpc

import (
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/judge/g"
	"github.com/open-falcon/judge/store"
	"time"
)

type Judge int

func (this *Judge) Ping(req model.NullRpcRequest, resp *model.SimpleRpcResponse) error {
	return nil
}

func (this *Judge) Send(items []*model.JudgeItem, resp *model.SimpleRpcResponse) error {
	remain := g.Config().Remain
	// 把当前时间的计算放在最外层，是为了减少获取时间时的系统调用开销
	now := time.Now().Unix()
	for _, item := range items {
		pk := item.PrimaryKey()
		store.HistoryBigMap[pk[0:2]].PushFrontAndMaintain(pk, item, remain, now)
	}
	return nil
}
