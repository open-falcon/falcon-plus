package rrdtool

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync/atomic"
	"time"

	pfc "github.com/niean/goperfcounter"
	"github.com/toolkits/consistent"

	cmodel "github.com/open-falcon/falcon-plus/common/model"
	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/graph/g"
	"github.com/open-falcon/falcon-plus/modules/graph/store"
)

const (
	_ = iota
	NET_TASK_M_SEND
	NET_TASK_M_QUERY
	NET_TASK_M_PULL
)

type Net_task_t struct {
	Method int
	Key    string
	Done   chan error
	Args   interface{}
	Reply  interface{}
}

const (
	FETCH_S_SUCCESS = iota
	FETCH_S_ERR
	FETCH_S_ISNOTEXIST
	SEND_S_SUCCESS
	SEND_S_ERR
	QUERY_S_SUCCESS
	QUERY_S_ERR
	CONN_S_ERR
	CONN_S_DIAL
	STAT_SIZE
)

var (
	Consistent       *consistent.Consistent
	Net_task_ch      map[string]chan *Net_task_t
	clients          map[string][]*rpc.Client
	flushrrd_timeout int32
	stat_cnt         [STAT_SIZE]uint64
)

func init() {
	Consistent = consistent.New()
	Net_task_ch = make(map[string]chan *Net_task_t)
	clients = make(map[string][]*rpc.Client)
}

func GetCounter() (ret string) {
	return fmt.Sprintf("FETCH_S_SUCCESS[%d] FETCH_S_ERR[%d] FETCH_S_ISNOTEXIST[%d] SEND_S_SUCCESS[%d] SEND_S_ERR[%d] QUERY_S_SUCCESS[%d] QUERY_S_ERR[%d] CONN_S_ERR[%d] CONN_S_DIAL[%d]",
		atomic.LoadUint64(&stat_cnt[FETCH_S_SUCCESS]),
		atomic.LoadUint64(&stat_cnt[FETCH_S_ERR]),
		atomic.LoadUint64(&stat_cnt[FETCH_S_ISNOTEXIST]),
		atomic.LoadUint64(&stat_cnt[SEND_S_SUCCESS]),
		atomic.LoadUint64(&stat_cnt[SEND_S_ERR]),
		atomic.LoadUint64(&stat_cnt[QUERY_S_SUCCESS]),
		atomic.LoadUint64(&stat_cnt[QUERY_S_ERR]),
		atomic.LoadUint64(&stat_cnt[CONN_S_ERR]),
		atomic.LoadUint64(&stat_cnt[CONN_S_DIAL]))
}

func dial(address string, timeout time.Duration) (*rpc.Client, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return rpc.NewClient(conn), err
}

func migrate_start(cfg *g.GlobalConfig) {
	var err error
	var i int
	if cfg.Migrate.Enabled {
		Consistent.NumberOfReplicas = cfg.Migrate.Replicas

		nodes := cutils.KeysOfMap(cfg.Migrate.Cluster)
		for _, node := range nodes {
			addr := cfg.Migrate.Cluster[node]
			Consistent.Add(node)
			Net_task_ch[node] = make(chan *Net_task_t, 16)
			clients[node] = make([]*rpc.Client, cfg.Migrate.Concurrency)

			for i = 0; i < cfg.Migrate.Concurrency; i++ {
				if clients[node][i], err = dial(addr, time.Second); err != nil {
					log.Fatalf("node:%s addr:%s err:%s\n", node, addr, err)
				}
				go net_task_worker(i, Net_task_ch[node], &clients[node][i], addr)
			}
		}
	}
}

func net_task_worker(idx int, ch chan *Net_task_t, client **rpc.Client, addr string) {
	var err error
	for {
		select {
		case task := <-ch:
			if task.Method == NET_TASK_M_SEND {
				if err = send_data(client, task.Key, addr); err != nil {
					pfc.Meter("migrate.send.err", 1)
					atomic.AddUint64(&stat_cnt[SEND_S_ERR], 1)
				} else {
					pfc.Meter("migrate.send.ok", 1)
					atomic.AddUint64(&stat_cnt[SEND_S_SUCCESS], 1)
				}
			} else if task.Method == NET_TASK_M_QUERY {
				if err = query_data(client, addr, task.Args, task.Reply); err != nil {
					pfc.Meter("migrate.query.err", 1)
					atomic.AddUint64(&stat_cnt[QUERY_S_ERR], 1)
				} else {
					pfc.Meter("migrate.query.ok", 1)
					atomic.AddUint64(&stat_cnt[QUERY_S_SUCCESS], 1)
				}
			} else if task.Method == NET_TASK_M_PULL {
				if atomic.LoadInt32(&flushrrd_timeout) != 0 {
					// hope this more faster than fetch_rrd
					if err = send_data(client, task.Key, addr); err != nil {
						pfc.Meter("migrate.sendbusy.err", 1)
						atomic.AddUint64(&stat_cnt[SEND_S_ERR], 1)
					} else {
						pfc.Meter("migrate.sendbusy.ok", 1)
						atomic.AddUint64(&stat_cnt[SEND_S_SUCCESS], 1)
					}
				} else {
					if err = fetch_rrd(client, task.Key, addr); err != nil {
						if os.IsNotExist(err) {
							pfc.Meter("migrate.scprrd.null", 1)
							//文件不存在时，直接将缓存数据刷入本地
							atomic.AddUint64(&stat_cnt[FETCH_S_ISNOTEXIST], 1)
							store.GraphItems.SetFlag(task.Key, 0)
							CommitByKey(task.Key)
						} else {
							pfc.Meter("migrate.scprrd.err", 1)
							//warning:其他异常情况，缓存数据会堆积
							atomic.AddUint64(&stat_cnt[FETCH_S_ERR], 1)
						}
					} else {
						pfc.Meter("migrate.scprrd.ok", 1)
						atomic.AddUint64(&stat_cnt[FETCH_S_SUCCESS], 1)
					}
				}
			} else {
				err = errors.New("error net task method")
			}
			if task.Done != nil {
				task.Done <- err
			}
		}
	}
}

// TODO addr to node
func reconnection(client **rpc.Client, addr string) {
	pfc.Meter("migrate.reconnection."+addr, 1)

	var err error

	atomic.AddUint64(&stat_cnt[CONN_S_ERR], 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = dial(addr, time.Second)
	atomic.AddUint64(&stat_cnt[CONN_S_DIAL], 1)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = dial(addr, time.Second)
		atomic.AddUint64(&stat_cnt[CONN_S_DIAL], 1)
	}
}

func query_data(client **rpc.Client, addr string,
	args interface{}, resp interface{}) error {
	var (
		err error
		i   int
	)

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.Query", args, resp,
			time.Duration(g.Config().CallTimeout)*time.Millisecond)

		if err == nil {
			break
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
	return err
}

func send_data(client **rpc.Client, key string, addr string) error {
	var (
		err  error
		flag uint32
		resp *cmodel.SimpleRpcResponse
		i    int
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
	resp = &cmodel.SimpleRpcResponse{}

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.Send", items, resp,
			time.Duration(cfg.CallTimeout)*time.Millisecond)

		if err == nil {
			goto out
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
	// err
	store.GraphItems.PushAll(key, items)
	//flag |= g.GRAPH_F_ERR
out:
	flag &= ^g.GRAPH_F_SENDING
	store.GraphItems.SetFlag(key, flag)
	return err

}

func fetch_rrd(client **rpc.Client, key string, addr string) error {
	var (
		err      error
		flag     uint32
		md5      string
		dsType   string
		filename string
		step, i  int
		rrdfile  g.File
	)

	cfg := g.Config()

	if flag, err = store.GraphItems.GetFlag(key); err != nil {
		return err
	}

	store.GraphItems.SetFlag(key, flag|g.GRAPH_F_FETCHING)

	md5, dsType, step, _ = g.SplitRrdCacheKey(key)
	filename = g.RrdFileName(cfg.RRD.Storage, md5, dsType, step)

	for i = 0; i < 3; i++ {
		err = rpc_call(*client, "Graph.GetRrd", key, &rrdfile,
			time.Duration(cfg.CallTimeout)*time.Millisecond)

		if err == nil {
			done := make(chan error, 1)
			io_task_chan <- &io_task_t{
				method: IO_TASK_M_WRITE,
				args: &g.File{
					Filename: filename,
					Body:     rrdfile.Body[:],
				},
				done: done,
			}
			if err = <-done; err != nil {
				goto out
			} else {
				flag &= ^g.GRAPH_F_MISS
				goto out
			}
		} else {
			log.Println(err)
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr)
		}
	}
out:
	flag &= ^g.GRAPH_F_FETCHING
	store.GraphItems.SetFlag(key, flag)
	return err
}

func rpc_call(client *rpc.Client, method string, args interface{},
	reply interface{}, timeout time.Duration) error {
	done := make(chan *rpc.Call, 1)
	client.Go(method, args, reply, done)
	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[rpc]")
	case call := <-done:
		if call.Error == nil {
			return nil
		} else {
			return call.Error
		}
	}
}
