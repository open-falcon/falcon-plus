package rpc

import (
	"strconv"
	"time"

	pfc "github.com/niean/goperfcounter"
	cmodel "github.com/open-falcon/common/model"
	cutils "github.com/open-falcon/common/utils"

	"github.com/open-falcon/gateway/g"
	"github.com/open-falcon/gateway/sender"
)

type Transfer int

func (this *Transfer) Ping(req cmodel.NullRpcRequest, resp *cmodel.SimpleRpcResponse) error {
	return nil
}

func (t *Transfer) Update(args []*cmodel.MetricValue, reply *g.TransferResp) error {
	return RecvMetricValues(args, reply, "rpc")
}

// process new metric values
func RecvMetricValues(args []*cmodel.MetricValue, reply *g.TransferResp, from string) error {
	start := time.Now()
	reply.ErrInvalid = 0

	items := []*cmodel.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.ErrInvalid += 1
			continue
		}

		// 历史遗留问题.
		// 老版本agent上报的metric=kernel.hostname的数据,其取值为string类型,现在已经不支持了;所以,这里硬编码过滤掉
		if v.Metric == "kernel.hostname" {
			reply.ErrInvalid += 1
			continue
		}

		if v.Metric == "" || v.Endpoint == "" {
			reply.ErrInvalid += 1
			continue
		}

		if v.Type != g.COUNTER && v.Type != g.GAUGE && v.Type != g.DERIVE {
			reply.ErrInvalid += 1
			continue
		}

		if v.Value == "" {
			reply.ErrInvalid += 1
			continue
		}

		if v.Step <= 0 {
			reply.ErrInvalid += 1
			continue
		}

		if len(v.Metric)+len(v.Tags) > 510 {
			reply.ErrInvalid += 1
			continue
		}

		errtags, tags := cutils.SplitTagsString(v.Tags)
		if errtags != nil {
			reply.ErrInvalid += 1
			continue
		}

		// TODO 呵呵,这里需要再优雅一点
		now := start.Unix()
		if v.Timestamp <= 0 || v.Timestamp > now*2 {
			v.Timestamp = now
		}

		fv := &cmodel.MetaData{
			Metric:      v.Metric,
			Endpoint:    v.Endpoint,
			Timestamp:   v.Timestamp,
			Step:        v.Step,
			CounterType: v.Type,
			Tags:        tags, //TODO tags键值对的个数,要做一下限制
		}

		valid := true
		var vv float64
		var err error

		switch cv := v.Value.(type) {
		case string:
			vv, err = strconv.ParseFloat(cv, 64)
			if err != nil {
				valid = false
			}
		case float64:
			vv = cv
		case int64:
			vv = float64(cv)
		default:
			valid = false
		}

		if !valid {
			reply.ErrInvalid += 1
			continue
		}

		fv.Value = vv
		items = append(items, fv)
	}

	// statistics
	cnt := int64(len(items))
	pfc.Meter("Recv", cnt)
	if from == "rpc" {
		pfc.Meter("RpcRecv", cnt)
	} else if from == "http" {
		pfc.Meter("HttpRecv", cnt)
	}

	cfg := g.Config()
	if cfg.Transfer.Enabled {
		sender.Push2SendQueue(items)
	}

	reply.Msg = "ok"
	reply.Total = len(args)
	reply.Latency = (time.Now().UnixNano() - start.UnixNano()) / 1000000

	return nil
}
