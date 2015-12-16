package rrdtool

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync/atomic"
	"time"

	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/graph/g"
	"github.com/open-falcon/graph/store"
)

func task_worker(idx int) {
	for {
		select {
		case key := <-task_key_ch:
			if atomic.LoadInt32(&flushrrd_timeout) != 0 {
				// hope this more faster than fetch_rrd
				send_data(key)
			} else {
				fetch_rrd(key)
			}
		}
	}
}

func send_data(key string) error {
	var (
		err    error
		flag   uint32
		node   string
		addr   string
		client *rpc.Client
		resp   *cmodel.SimpleRpcResponse
	)

	//remote
	if flag, err = store.GraphItems.GetFlag(key); err != nil {
		return err
	}
	cfg := g.Config()

	store.GraphItems.SetFlag(key, flag|g.GRAPH_F_SENDING)

	items := store.GraphItems.PopAll(key)
	items_size := len(items)
	if items_size == 0 {
		goto out
	}

	node, _ = Consistent.Get(items[0].PrimaryKey())
	client = Client[node]
	resp = &cmodel.SimpleRpcResponse{}

	err = Jsonrpc_call(client, "Graph.Send", items, resp,
		time.Duration(cfg.CallTimeout)*time.Millisecond)

	// reconnection
	if err != nil {
		store.GraphItems.PushAll(key, items)

		conn.Lock()
		client.Close()
		addr = cfg.Migrate.Cluster[node]
		client, err = jsonrpc.Dial("tcp", addr)
		conn.Unlock()

		for err != nil {
			//danger!! block routine
			time.Sleep(time.Millisecond * 500)
			conn.Lock()
			client, err = jsonrpc.Dial("tcp", addr)
			conn.Unlock()
		}
		goto err_out
	}
	goto out

err_out:
	flag |= g.GRAPH_F_ERR
out:
	flag &= ^g.GRAPH_F_SENDING
	store.GraphItems.SetFlag(key, flag)
	return err

}

//func fetch_rrd(client *rpc.Client, queue *store.SafeLinkedList, node, addr string) {
func fetch_rrd(key string) error {
	var (
		err      error
		flag     uint32
		md5      string
		dsType   string
		filename string
		step     int
		rrdfile  g.File64
		ctx      []byte
		node     string
		addr     string
		client   *rpc.Client
	)

	cfg := g.Config()

	if flag, err = store.GraphItems.GetFlag(key); err != nil {
		return err
	}

	store.GraphItems.SetFlag(key, flag|g.GRAPH_F_FETCHING)

	md5, dsType, step, _ = g.SplitRrdCacheKey(key)
	filename = g.RrdFileName(cfg.RRD.Storage, md5, dsType, step)

	items := store.GraphItems.PopAll(key)
	items_size := len(items)
	if items_size == 0 {
		// impossible
		goto out
	}

	node, _ = Consistent.Get(items[0].PrimaryKey())
	client = Client[node]

	err = Jsonrpc_call(client, "Graph.GetRrd", key, &rrdfile,
		time.Duration(cfg.CallTimeout)*time.Millisecond)

	// reconnection
	if err != nil {
		store.GraphItems.PushAll(key, items)

		client.Close()
		addr = cfg.Migrate.Cluster[node]
		client, err = jsonrpc.Dial("tcp", addr)
		for err != nil {
			//danger!! block routine
			time.Sleep(time.Millisecond * 500)
			client, err = jsonrpc.Dial("tcp", addr)
		}
		goto err_out
	}

	if ctx, err = base64.StdEncoding.DecodeString(rrdfile.Body64); err != nil {
		store.GraphItems.PushAll(key, items)
		goto err_out
	} else {
		if err = ioutil.WriteFile(filename, ctx, 0644); err != nil {
			store.GraphItems.PushAll(key, items)
			goto err_out
		} else {
			flag &= ^g.GRAPH_F_MISS
			Flush(filename, items)
			goto out
		}
	}
	//noneed
	goto out

err_out:
	flag |= g.GRAPH_F_ERR
out:
	flag &= ^g.GRAPH_F_FETCHING
	store.GraphItems.SetFlag(key, flag)
	return err
}

func Jsonrpc_call(client *rpc.Client, method string, args interface{},
	reply interface{}, timeout time.Duration) error {
	done := make(chan *rpc.Call, 1)
	client.Go(method, args, reply, done)
	select {
	case <-time.After(timeout):
		return errors.New("timeout")
	case call := <-done:
		if call.Error == nil {
			return nil
		} else {
			return call.Error
		}
	}
}
