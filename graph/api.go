package graph

import (
	"errors"
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/query/g"
	"github.com/toolkits/logger"
	"github.com/toolkits/rpool/conn_pool"
	"math"
	"sync"
	"time"
)

func Info(endpoint, counter string) (r *model.GraphFullyInfo, err error) {
	pool, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}

	rpc_conn := conn.(RpcConn)
	if rpc_conn.cli == nil {
		pool.CloseClean(conn)
		return nil, errors.New("nil rpc conn")
	}

	type ChResult struct {
		Err  error
		Resp *model.GraphInfoResp
	}
	ch := make(chan *ChResult, 1)
	go func() {
		param := model.GraphInfoParam{
			Endpoint: endpoint,
			Counter:  counter,
		}
		resp := &model.GraphInfoResp{}
		err := rpc_conn.cli.Call("Graph.Info", param, resp)
		r := &ChResult{
			Err:  err,
			Resp: resp,
		}
		ch <- r
	}()

	cfg := g.Config().Graph
	select {
	case r := <-ch:
		if r.Err != nil {
			pool.CloseClean(conn)
			return nil, r.Err
		} else {
			pool.Release(conn)
			logger.Trace("graph.info resp: %v, addr: %v", r.Resp, pool.Name)
			fullyInfo := model.GraphFullyInfo{
				Endpoint:  endpoint,
				Counter:   counter,
				ConsolFun: r.Resp.ConsolFun,
				Step:      r.Resp.Step,
				Filename:  r.Resp.Filename,
				Addr:      pool.Name,
			}
			return &fullyInfo, nil
		}

	case <-time.After(time.Duration(cfg.Timeout) * time.Millisecond):
		pool.Release(conn)
		logger.Trace("graph.info timeout err: i/o timeout, addr: %v", pool.Name)
		return nil, errors.New("i/o timeout")
	}
}

func Last(endpoint, counter string) (r *model.GraphLastResp, err error) {
	pool, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}

	rpc_conn := conn.(RpcConn)
	if rpc_conn.cli == nil {
		pool.CloseClean(conn)
		return nil, errors.New("nil rpc conn")
	}

	type ChResult struct {
		Err  error
		Resp *model.GraphLastResp
	}
	ch := make(chan *ChResult, 1)
	go func() {
		param := model.GraphLastParam{
			Endpoint: endpoint,
			Counter:  counter,
		}
		resp := &model.GraphLastResp{}
		err := rpc_conn.cli.Call("Graph.Last", param, resp)
		r := &ChResult{
			Err:  err,
			Resp: resp,
		}
		ch <- r
	}()

	cfg := g.Config().Graph
	select {
	case r := <-ch:
		if r.Err != nil {
			pool.CloseClean(conn)
			return nil, r.Err
		} else {
			pool.Release(conn)
			return r.Resp, nil
		}
	case <-time.After(time.Duration(cfg.Timeout) * time.Millisecond):
		pool.Release(conn)
		return nil, errors.New("i/o timeout")
	}
}

func QueryOne(start, end int64, cf, endpoint, counter string) (r *model.GraphQueryResponse, err error) {
	pool, err := selectPool(endpoint, counter)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Get()
	if err != nil {
		return nil, err
	}

	rpc_conn := conn.(RpcConn)
	if rpc_conn.cli == nil {
		pool.CloseClean(conn)
		return nil, errors.New("nil rpc conn")
	}

	type ChResult struct {
		Err  error
		Resp *model.GraphQueryResponse
	}
	ch := make(chan *ChResult, 1)
	go func() {
		param := model.GraphQueryParam{
			Start:     start,
			End:       end,
			ConsolFun: cf,
			Endpoint:  endpoint,
			Counter:   counter,
		}
		resp := &model.GraphQueryResponse{}
		err := rpc_conn.cli.Call("Graph.Query", param, resp)
		r := &ChResult{
			Err:  err,
			Resp: resp,
		}
		ch <- r
	}()

	cfg := g.Config().Graph
	select {
	case r := <-ch:
		if r.Err != nil {
			pool.CloseClean(conn)
			return nil, r.Err
		} else {
			pool.Release(conn)
			logger.Trace("graph: query graph resp: %v, addr: %v", r.Resp, pool.Name)

			fixedResp := &model.GraphQueryResponse{
				Endpoint: r.Resp.Endpoint,
				Counter:  r.Resp.Counter,
				DsType:   r.Resp.DsType,
				Step:     r.Resp.Step,
			}
			size := len(r.Resp.Values)

			//NOTICE:最后一个点是坏点，过滤点，可能是rrdtool的bug
			if size < 1 {
				return fixedResp, nil
			} else {
				dsType := r.Resp.DsType
				fixedValues := []*model.RRDData{}
				for _, v := range r.Resp.Values[0:size] {
					if v == nil {
						continue
					}
					if v.Timestamp < start || v.Timestamp > end {
						continue
					}
					//FIXME: 查询数据的时候，把所有的负值都过滤掉，因为transfer之前在设置最小值的时候为U
					if (dsType == "DERIVE" || dsType == "COUNTER") && v.Value < 0 {
						fixedValues = append(fixedValues, &model.RRDData{
							Timestamp: v.Timestamp,
							Value:     model.JsonFloat(math.NaN()),
						})
					} else {
						fixedValues = append(fixedValues, v)
					}
				}
				fixedResp.Values = fixedValues
				return fixedResp, nil
			}
		}

	case <-time.After(time.Duration(cfg.Timeout) * time.Millisecond):
		pool.Release(conn)
		logger.Trace("query graph timeout err: i/o timeout, addr: %v", pool.Name)
		return nil, errors.New("i/o timeout")
	}

}

// metrics: [["endpoint", "counter"], []]
func QueryMulti(start, end int64, cf string, metrics [][]string) []*model.GraphQueryResponse {
	result := make([]*model.GraphQueryResponse, 0, len(metrics))

	var wg sync.WaitGroup
	lock := new(sync.Mutex)

	for _, pair := range metrics {
		if len(pair) != 2 {
			continue
		}
		endpoint := pair[0]
		counter := pair[1]

		wg.Add(1)
		go func(endpoint, counter string) {
			defer wg.Done()
			r, err := QueryOne(start, end, cf, endpoint, counter)
			if err != nil {
				logger.Error("query one from graph fail: %v", err)
				return
			}

			lock.Lock()
			defer lock.Unlock()
			result = append(result, r)

		}(endpoint, counter)
	}

	wg.Wait()

	return result
}

func selectPool(endpoint, counter string) (*conn_pool.ConnPool, error) {
	pkey := fmt.Sprintf("%s/%s", endpoint, counter)

	hash_node, err := backend.LocateRing(pkey)
	if err != nil {
		return nil, err
	}

	pools, err := backend.GetConnPoolsByName(hash_node)
	if err != nil {
		return nil, err
	}

	if len(pools) == 0 {
		return nil, errors.New("empty_conn_pool")
	}

	return pools[0], nil
}
