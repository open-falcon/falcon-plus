// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl"
)

type Kv struct {
	Key     string
	Section string
	Value   string
}

var (
	etcdMap = map[string]map[string]string{
		"graph": map[string]string{
			falcon.C_DEBUG:                "/open-falcon/graph/config/debug",
			falcon.C_HTTP_ENABLE:          "/open-falcon/graph/config/http/enabled",
			falcon.C_HTTP_ADDR:            "/open-falcon/graph/config/http/listen",
			falcon.C_RPC_ENABLE:           "/open-falcon/graph/config/rpc/enabled",
			falcon.C_RPC_ADDR:             "/open-falcon/graph/config/rpc/listen",
			falcon.C_RRD_STORAGE:          "/open-falcon/graph/config/rrd/storage",
			falcon.C_DSN:                  "/open-falcon/graph/config/db/dsn",
			falcon.C_DB_MAX_IDLE:          "/open-falcon/graph/config/db/maxIdle",
			falcon.C_CALL_TIMEOUT:         "/open-falcon/graph/config/callTimeout",
			falcon.C_MIGRATE_ENABLE:       "/open-falcon/graph/config/migrate/enabled",
			falcon.C_MIGRATE_CONCURRENCY:  "/open-falcon/graph/config/migrate/concurrency",
			falcon.C_MIGRATE_REPLICAS:     "/open-falcon/graph/config/migrate/replicas",
			falcon.C_MIGRATE_CLUSTER:      "/open-falcon/graph/config/migrate/cluster",
			falcon.C_MIGRATE_NEW_ENDPOINT: "/open-falcon/graph/config/migrate/newEndpoint",
			falcon.C_LEASE_TTL:            "/open-falcon/graph/config/leasettl",
			falcon.C_GRPC_ENABLE:          "/open-falcon/graph/config/grpc/enabled",
			falcon.C_GRPC_ADDR:            "/open-falcon/graph/config/grpc/listen",
		},
		"transfer": map[string]string{
			falcon.C_DEBUG:             "/open-falcon/transfer/config/debug",
			falcon.C_MINSTEP:           "/open-falcon/transfer/config/minStep",
			falcon.C_HTTP_ENABLE:       "/open-falcon/transfer/config/http/enabled",
			falcon.C_HTTP_ADDR:         "/open-falcon/transfer/config/http/listen",
			falcon.C_RPC_ENABLE:        "/open-falcon/transfer/config/rpc/enabled",
			falcon.C_RPC_ADDR:          "/open-falcon/transfer/config/rpc/listen",
			falcon.C_SOCKET_ENABLE:     "/open-falcon/transfer/config/socket/enabled",
			falcon.C_SOCKET_ADDR:       "/open-falcon/transfer/config/socket/listen",
			falcon.C_SOCKET_TIMEOUT:    "/open-falcon/transfer/config/socket/timeout",
			falcon.C_JUDGE_ENABLE:      "/open-falcon/transfer/config/judge/enabled",
			falcon.C_JUDGE_BATCH:       "/open-falcon/transfer/config/judge/batch",
			falcon.C_JUDGE_CONNTIMEOUT: "/open-falcon/transfer/config/judge/connTimeout",
			falcon.C_JUDGE_CALLTIMEOUT: "/open-falcon/transfer/config/judge/callTimeout",
			falcon.C_JUDGE_MAXCONNS:    "/open-falcon/transfer/config/judge/maxConns",
			falcon.C_JUDGE_MAXIDLE:     "/open-falcon/transfer/config/judge/maxIdle",
			falcon.C_JUDGE_REPLICAS:    "/open-falcon/transfer/config/judge/replicas",
			falcon.C_JUDGE_CLUSTER:     "/open-falcon/transfer/config/judge/cluster",
			falcon.C_GRAPH_ENABLE:      "/open-falcon/transfer/config/graph/enabled",
			falcon.C_GRAPH_BATCH:       "/open-falcon/transfer/config/graph/batch",
			falcon.C_GRAPH_CONNTIMEOUT: "/open-falcon/transfer/config/graph/connTimeout",
			falcon.C_GRAPH_CALLTIMEOUT: "/open-falcon/transfer/config/graph/callTimeout",
			falcon.C_GRAPH_MAXCONNS:    "/open-falcon/transfer/config/graph/maxConns",
			falcon.C_GRAPH_MAXIDLE:     "/open-falcon/transfer/config/graph/maxIdle",
			falcon.C_GRAPH_REPLICAS:    "/open-falcon/transfer/config/graph/replicas",
			falcon.C_GRAPH_CLUSTER:     "/open-falcon/transfer/config/graph/cluster",
			falcon.C_TSDB_ENABLE:       "/open-falcon/transfer/config/tsdb/enabled",
			falcon.C_TSDB_BATCH:        "/open-falcon/transfer/config/tsdb/batch",
			falcon.C_TSDB_CONNTIMEOUT:  "/open-falcon/transfer/config/tsdb/connTimeout",
			falcon.C_TSDB_CALLTIMEOUT:  "/open-falcon/transfer/config/tsdb/callTimeout",
			falcon.C_TSDB_MAXCONNS:     "/open-falcon/transfer/config/tsdb/maxConns",
			falcon.C_TSDB_MAXIDLE:      "/open-falcon/transfer/config/tsdb/maxIdle",
			falcon.C_TSDB_RETRY:        "/open-falcon/transfer/config/tsdb/retry",
			falcon.C_TSDB_ADDRESS:      "/open-falcon/transfer/config/tsdb/address",
			falcon.C_LEASE_TTL:         "/open-falcon/transfer/config/leasettl",
		},
	}
)

func prepareEtcdConfig() error {
	put := make(map[string]string)
	for module, kv := range etcdMap {
		ks := make(map[string]bool)
		for _, k := range kv {
			ks[k] = false
		}

		prefix := fmt.Sprintf("/open-falcon/%s/config/", module)
		resp, err := ctrl.EtcdCli.GetPrefix(prefix)
		if err != nil {
			return err
		}
		for _, v := range resp.Kvs {
			if _, ok := ks[string(v.Key)]; ok {
				ks[string(v.Key)] = true
			}
		}
		for k, exist := range ks {
			if !exist {
				put[k] = ""
			}
		}
	}
	return ctrl.EtcdCli.Puts(put)
}

func GetDbConfig(o orm.Ormer, module string) (ret map[string]string, err error) {
	var row Kv

	err = o.Raw("SELECT `section`, `key`, `value` FROM `kv` where "+
		"`section` = ? and `key` = 'config'", module).QueryRow(&row)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(row.Value), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (op *Operator) SetEtcdConfig(module string, conf map[string]string) error {
	ks, ok := etcdMap[module]
	if !ok {
		//skip miss hit
		return nil
	}
	ekv := make(map[string]string)
	for k, ek := range ks {
		beego.Debug(k, "->", ek, "=", conf[k])
		if conf[k] != "" {
			ekv[ek] = conf[k]
		}
	}
	return ctrl.EtcdCli.Puts(ekv)
}

func (op *Operator) SetDbConfig(module string, conf map[string]string) error {
	kv, _ := GetDbConfig(op.O, module)
	for k, v := range conf {
		if v != "" {
			kv[k] = v
		}
	}
	v, err := json.Marshal(kv)
	if err != nil {
		return err
	}
	s := string(v)
	_, err = op.O.Raw("INSERT INTO `kv`(`section`, `key`, `value`)"+
		" VALUES (?,'config',?) ON DUPLICATE KEY UPDATE `value`=?",
		module, s, s).Exec()

	return err
}

func (op *Operator) ConfigGet(module string) (interface{}, error) {
	var c *falcon.Configer

	switch module {
	case "ctrl":
		c = &ctrl.Configure.Ctrl
	case "agent":
		c = &ctrl.Configure.Agent
	case "loadbalance":
		c = &ctrl.Configure.Loadbalance
	case "backend":
		c = &ctrl.Configure.Backend
	case "graph": // for falcon-plus
		c = &ctrl.Configure.Graph
	case "transfer": // for falcon-plus
		c = &ctrl.Configure.Transfer
	default:
		return nil, ErrNoModule
	}

	conf, err := GetDbConfig(op.O, module)
	if err == nil {
		c.Set(falcon.APP_CONF_DB, conf)
	}
	return c.Get(), nil
}

func (op *Operator) ConfigSet(module string, conf map[string]string) error {
	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		if module == "graph" {
			// disable set expansion from this interface
			delete(conf, falcon.C_MIGRATE_ENABLE)
			delete(conf, falcon.C_MIGRATE_NEW_ENDPOINT)
		}
		err := op.SetEtcdConfig(module, conf)
		if err != nil {
			return err
		}
		return op.SetDbConfig(module, conf)
	default:
		return ErrNoModule
	}
}

func (op *Operator) OnlineGet(module string) ([]KeyValue, error) {

	switch module {
	case "ctrl", "agent", "lb", "backend", "graph", "transfer":
		prefix := fmt.Sprintf("/open-falcon/%s/online/", module)
		resp, err := ctrl.EtcdCli.GetPrefix(prefix)
		if err != nil {
			return nil, err
		}

		ret := make([]KeyValue, len(resp.Kvs))

		for i, kv := range resp.Kvs {
			ret[i] = KeyValue{
				Key:   string(kv.Key[len(prefix):]),
				Value: string(kv.Value),
			}
		}

		return ret, nil
	default:
		return nil, ErrNoModule
	}

}

type ExpansionStatus struct {
	Migrating    bool   `json:"migrating"`
	GraphCluster string `json:"graph_cluster"`
	NewEndpoint  string `json:"new_endpoint"`
}

// expansion
func (op *Operator) ExpansionGet(module string) (ret *ExpansionStatus, err error) {
	var enable string
	ret = &ExpansionStatus{}

	if module != "graph" {
		return nil, ErrUnsupported
	}

	if enable, err = ctrl.EtcdCli.Get(etcdMap["graph"][falcon.C_MIGRATE_ENABLE]); err == nil && enable == "true" {
		ret.Migrating = true
	}

	ret.GraphCluster, _ = ctrl.EtcdCli.Get(etcdMap["transfer"][falcon.C_GRAPH_CLUSTER])
	ret.NewEndpoint, _ = ctrl.EtcdCli.Get(etcdMap["graph"][falcon.C_MIGRATE_NEW_ENDPOINT])

	return ret, nil
}

func (op *Operator) ExpansionBegin(module string, newEndpoint string) error {
	if module != "graph" {
		return ErrUnsupported
	}

	// get
	transfer, err := GetDbConfig(op.O, "transfer")
	if err != nil {
		return err
	}

	replicas, ok := transfer[falcon.C_GRAPH_REPLICAS]
	if !ok {
		return errors.New("can not get transfer->graph->replicas config")
	}

	online := make(map[string]string)
	_online, _ := op.OnlineGet("graph")
	for _, kv := range _online {
		online[kv.Key] = kv.Value
	}

	old_cluster := make(map[string]string)
	new_cluster := make(map[string]string)
	_old_cluster, ok := transfer[falcon.C_GRAPH_CLUSTER]
	if !ok {
		return errors.New("can not get transfer->graph->cluster config")
	}
	for _, v := range strings.Split(_old_cluster, ";") {
		// alias=hostname
		// s[0]: alias
		// s[1]: hostname
		s := strings.Split(v, "=")
		old_cluster[s[1]] = s[0]
		new_cluster[s[0]] = s[1]
	}

	for _, v := range strings.Split(newEndpoint, ";") {
		s := strings.Split(v, "=")
		if len(s) != 2 {
			return errors.New("endpoint format error " + v)
		}
		h := strings.Split(s[1], ":")
		if _, ok := online[h[0]]; !ok {
			return errors.New(fmt.Sprintf("endpoint(%s) is not online", h[0]))
		}
		new_cluster[s[0]] = s[1]
	}

	_new_cluster := ""
	for k, v := range new_cluster {
		_new_cluster += fmt.Sprintf("%s=%s;", k, v)
	}
	if len(_new_cluster) > 0 {
		_new_cluster = _new_cluster[0 : len(_new_cluster)-1]
	}

	// set
	op.SetDbConfig("graph", map[string]string{
		falcon.C_MIGRATE_NEW_ENDPOINT: newEndpoint,
		falcon.C_MIGRATE_ENABLE:       "true",
		falcon.C_MIGRATE_REPLICAS:     replicas,
		falcon.C_MIGRATE_CLUSTER:      _old_cluster,
	})
	op.SetDbConfig("transfer", map[string]string{
		falcon.C_GRAPH_CLUSTER: _new_cluster,
	})

	ekv := make(map[string]string)
	ekv[etcdMap["graph"][falcon.C_MIGRATE_NEW_ENDPOINT]] = newEndpoint
	ekv[etcdMap["graph"][falcon.C_MIGRATE_ENABLE]] = "true"
	ekv[etcdMap["graph"][falcon.C_MIGRATE_REPLICAS]] = replicas
	ekv[etcdMap["graph"][falcon.C_MIGRATE_CLUSTER]] = _old_cluster
	ekv[etcdMap["transfer"][falcon.C_GRAPH_CLUSTER]] = _new_cluster

	return ctrl.EtcdCli.Puts(ekv)

}

func (op *Operator) ExpansionFinish(module string) error {
	if module != "graph" {
		return ErrUnsupported
	}

	op.SetDbConfig("graph", map[string]string{
		falcon.C_MIGRATE_ENABLE:       "false",
		falcon.C_MIGRATE_NEW_ENDPOINT: " ",
	})

	ekv := make(map[string]string)
	ekv[etcdMap["graph"][falcon.C_MIGRATE_ENABLE]] = "false"
	ekv[etcdMap["graph"][falcon.C_MIGRATE_NEW_ENDPOINT]] = " "
	return ctrl.EtcdCli.Puts(ekv)
}
