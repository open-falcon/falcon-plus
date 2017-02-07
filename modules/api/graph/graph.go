package graph

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"

	backend "github.com/open-falcon/falcon-plus/common/backend_pool"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/spf13/viper"
	connp "github.com/toolkits/conn_pool"
	rpcpool "github.com/toolkits/conn_pool/rpc_conn_pool"
	rings "github.com/toolkits/consistent/rings"
	nset "github.com/toolkits/container/set"
)

// 连接池
// node_address -> connection_pool
var (
	GraphConnPools *backend.SafeRpcConnPools
	clusterMap     map[string]string
	gcluster       []string
	connTimeout    int32
	callTimeout    int32
)

// 服务节点的一致性哈希环
// pk -> node
var (
	GraphNodeRing *rings.ConsistentHashNodeRing
)

func Start(addrs map[string]string) {
	clusterMap = addrs
	connTimeout = int32(viper.GetInt("graphs.conn_timeout"))
	callTimeout = int32(viper.GetInt("graphs.call_timeout"))
	for c := range clusterMap {
		gcluster = append(gcluster, c)
	}
	defer func() {
		if r := recover(); r != nil {
			log.Error("graph got painc")
			log.Error(fmt.Sprintf("%s", r))
			Start(clusterMap)
		}
	}()
	initNodeRings(clusterMap)
	initConnPools(clusterMap)
	log.Println("graph.Start ok")
}

func GenQParam(endpoint string, counter string, consolFun string, stime int64, etime int64, step int) cmodel.GraphQueryParam {
	return cmodel.GraphQueryParam{
		Start:     stime,
		End:       etime,
		ConsolFun: consolFun,
		Endpoint:  endpoint,
		Counter:   counter,
		Step:      step,
	}
}
func QueryOne(para cmodel.GraphQueryParam) (resp *cmodel.GraphQueryResponse, err error) {
	start, end := para.Start, para.End
	endpoint, counter := para.Endpoint, para.Counter
	resp = &cmodel.GraphQueryResponse{}
	pool, addr, err := selectPool(endpoint, counter)
	if err != nil {
		return resp, err
	}

	conn, err := pool.Fetch()
	if err != nil {
		return resp, err
	}

	rpcConn := conn.(*rpcpool.RpcClient)
	if rpcConn.Closed() {
		pool.ForceClose(conn)
		return resp, errors.New("conn closed")
	}

	type ChResult struct {
		Err  error
		Resp *cmodel.GraphQueryResponse
	}

	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphQueryResponse{}
		err := rpcConn.Call("Graph.Query", para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		pool.ForceClose(conn)
		return nil, fmt.Errorf("%s, call timeout. proc: %s", addr, pool.Proc())
	case r := <-ch:
		if r.Err != nil {
			pool.ForceClose(conn)
			return r.Resp, fmt.Errorf("%s, call failed, err %v. proc: %s", addr, r.Err, pool.Proc())
		} else {
			pool.Release(conn)

			if len(r.Resp.Values) < 1 {
				r.Resp.Values = []*cmodel.RRDData{}
				return r.Resp, nil
			}

			// TODO query不该做这些事情, 说明graph没做好
			fixed := []*cmodel.RRDData{}
			for _, v := range r.Resp.Values {
				if v == nil || !(v.Timestamp >= start && v.Timestamp <= end) {
					continue
				}
				//FIXME: 查询数据的时候，把所有的负值都过滤掉，因为transfer之前在设置最小值的时候为U
				if (r.Resp.DsType == "DERIVE" || r.Resp.DsType == "COUNTER") && v.Value < 0 {
					fixed = append(fixed, &cmodel.RRDData{Timestamp: v.Timestamp, Value: cmodel.JsonFloat(math.NaN())})
				} else {
					fixed = append(fixed, v)
				}
			}
			r.Resp.Values = fixed
		}
		return r.Resp, nil
	}
}

func Info(para cmodel.GraphInfoParam) (resp *cmodel.GraphFullyInfo, err error) {
	endpoint, counter := para.Endpoint, para.Counter

	pool, addr, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Fetch()
	if err != nil {
		return nil, err
	}

	rpcConn := conn.(*rpcpool.RpcClient)
	if rpcConn.Closed() {
		pool.ForceClose(conn)
		return nil, errors.New("conn closed")
	}

	type ChResult struct {
		Err  error
		Resp *cmodel.GraphInfoResp
	}
	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphInfoResp{}
		err := rpcConn.Call("Graph.Info", para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		pool.ForceClose(conn)
		return nil, fmt.Errorf("%s, call timeout. proc: %s", addr, pool.Proc())
	case r := <-ch:
		if r.Err != nil {
			pool.ForceClose(conn)
			return nil, fmt.Errorf("%s, call failed, err %v. proc: %s", addr, r.Err, pool.Proc())
		} else {
			pool.Release(conn)
			fullyInfo := cmodel.GraphFullyInfo{
				Endpoint:  endpoint,
				Counter:   counter,
				ConsolFun: r.Resp.ConsolFun,
				Step:      r.Resp.Step,
				Filename:  r.Resp.Filename,
				Addr:      addr,
			}
			return &fullyInfo, nil
		}
	}
}

func Last(para cmodel.GraphLastParam) (r *cmodel.GraphLastResp, err error) {
	endpoint, counter := para.Endpoint, para.Counter

	pool, addr, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Fetch()
	if err != nil {
		return nil, err
	}

	rpcConn := conn.(*rpcpool.RpcClient)
	if rpcConn.Closed() {
		pool.ForceClose(conn)
		return nil, errors.New("conn closed")
	}

	type ChResult struct {
		Err  error
		Resp *cmodel.GraphLastResp
	}
	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphLastResp{}
		err := rpcConn.Call("Graph.Last", para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		pool.ForceClose(conn)
		return nil, fmt.Errorf("%s, call timeout. proc: %s", addr, pool.Proc())
	case r := <-ch:
		if r.Err != nil {
			pool.ForceClose(conn)
			return r.Resp, fmt.Errorf("%s, call failed, err %v. proc: %s", addr, r.Err, pool.Proc())
		} else {
			pool.Release(conn)
			return r.Resp, nil
		}
	}
}

func LastRaw(para cmodel.GraphLastParam) (r *cmodel.GraphLastResp, err error) {
	endpoint, counter := para.Endpoint, para.Counter

	pool, addr, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Fetch()
	if err != nil {
		return nil, err
	}

	rpcConn := conn.(*rpcpool.RpcClient)
	if rpcConn.Closed() {
		pool.ForceClose(conn)
		return nil, errors.New("conn closed")
	}

	type ChResult struct {
		Err  error
		Resp *cmodel.GraphLastResp
	}
	ch := make(chan *ChResult, 1)
	go func() {
		resp := &cmodel.GraphLastResp{}
		err := rpcConn.Call("Graph.LastRaw", para, resp)
		ch <- &ChResult{Err: err, Resp: resp}
	}()

	select {
	case <-time.After(time.Duration(callTimeout) * time.Millisecond):
		pool.ForceClose(conn)
		return nil, fmt.Errorf("%s, call timeout. proc: %s", addr, pool.Proc())
	case r := <-ch:
		if r.Err != nil {
			pool.ForceClose(conn)
			return r.Resp, fmt.Errorf("%s, call failed, err %v. proc: %s", addr, r.Err, pool.Proc())
		} else {
			pool.Release(conn)
			return r.Resp, nil
		}
	}
}

func selectPool(endpoint, counter string) (rpool *connp.ConnPool, raddr string, rerr error) {
	pkey := cutils.PK2(endpoint, counter)
	node, err := GraphNodeRing.GetNode(pkey)
	if err != nil {
		return nil, "", err
	}

	addr, found := clusterMap[node]
	if !found {
		return nil, "", errors.New("node not found")
	}

	pool, found := GraphConnPools.Get(addr)
	if !found {
		log.Errorf("pool :%v", pool)
		return nil, addr, errors.New("addr not found")
	}

	return pool, addr, nil
}

// internal functions
func initConnPools(clusterMap map[string]string) {

	// TODO 为了得到Slice,这里做的太复杂了
	graphInstances := nset.NewSafeSet()
	for _, address := range clusterMap {
		graphInstances.Add(address)
	}
	GraphConnPools = backend.CreateSafeRpcConnPools(
		int(viper.GetInt("graphs.max_conns")),
		int(viper.GetInt("graphs.max_idle")),
		int(connTimeout), int(callTimeout), graphInstances.ToSlice())
}

func initNodeRings(clusterMap map[string]string) {
	gcluster := cutils.KeysOfMap(clusterMap)
	GraphNodeRing = rings.NewConsistentHashNodesRing(
		int32(viper.GetInt("graphs.numberOfReplicas")),
		gcluster)
}

func Hosts() []string {
	f, _ := ioutil.ReadFile("hosts")
	splitLine := strings.Split(string(f), "\n")
	ss := []string{}

	for _, d := range splitLine {
		if strings.Contains(d, ",") {
			m := strings.Split(d, ",")
			ss = append(ss, m[1])
		} else {
			log.Debug(d)
		}
	}
	return ss
}
