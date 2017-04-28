// Copyright 2017 Xiaomi, Inc.
// Copyright 2014 beego Author. All Rights Reserved.
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
package falcon

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	APP_CONF_DEFAULT = iota
	APP_CONF_DB
	APP_CONF_FILE
	APP_CONF_SIZE
)

const (
	C_ETCD_ENDPOINTS          = "etcdendpoints"
	C_ETCD_USERNAME           = "etcdusername"
	C_ETCD_PASSWORD           = "etcdpassword"
	C_ETCD_CERTFILE           = "certfile"
	C_ETCD_KEYFILE            = "keyfile"
	C_ETCD_CAFILE             = "cafile"
	C_LEASE_KEY               = "leasekey"
	C_LEASE_VALUE             = "leasevalue"
	C_LEASE_TTL               = "leasettl"
	C_RUN_MODE                = "runmode"
	C_ENABLE_DOCS             = "enabledocs"
	C_SEESION_GC_MAX_LIFETIME = "sessiongcmaxlifetime"
	C_SESSION_COOKIE_LIFETIME = "sessioncookielifetime"
	C_AUTH_MODULE             = "authmodule"
	C_CACHE_MODULE            = "cachemodule"
	C_LDAP_ADDR               = "ldapaddr"
	C_LDAP_BASE_DN            = "ldapbasedn"
	C_LDAP_BIND_DN            = "ldapbinddn"
	C_LDAP_BIND_PWD           = "ldapbindpwd"
	C_LDAP_FILTER             = "ldapfilter"
	C_MISSO_REDIRECT_URL      = "missoredirecturl"
	C_GITHUB_CLIENT_ID        = "githubclientid"
	C_GITHUB_CLIENT_SECRET    = "githubclientsecret"
	C_GITHUB_REDIRECT_URL     = "githubredirecturl"
	C_GOOGLE_CLIENT_ID        = "googleclientid"
	C_GOOGLE_CLIENT_SECRET    = "googleclientsecret"
	C_GOOGLE_REDIRECT_URL     = "googleredirecturl"
	C_TAG_SCHEMA              = "tagschema"
	C_CONN_TIMEOUT            = "conntimeout"
	C_CALL_TIMEOUT            = "calltimeout"
	C_WORKER_PROCESSES        = "workerprocesses"
	C_HTTP_ENABLE             = "httpenable"
	C_HTTP_ADDR               = "httpaddr"
	C_RPC_ENABLE              = "rpcenable"
	C_RPC_ADDR                = "rpcaddr"
	C_GRPC_ENABLE             = "grpcenable"
	C_GRPC_ADDR               = "grpcaddr"
	C_INTERVAL                = "interval"
	C_IFACE_PREFIX            = "ifaceprefix"
	C_PAYLOADSIZE             = "payloadsize"
	C_DSN                     = "dsn"
	C_DB_MAX_IDLE             = "dbmaxidle"
	C_DB_MAX_CONN             = "dbmaxconn"
	C_HDISK                   = "hdisk"
	C_UPSTREAM                = "upstream"
	C_IDX                     = "idx"
	C_IDXINTERVAL             = "idxinterval"
	C_IDXFULLINTERVAL         = "idxfullinterval"
	C_SHMMAGIC                = "shmmagic"
	C_SHMKEY                  = "shmkey"
	C_SHMSIZE                 = "shmsize"
	C_EMU_ENABLE              = "emuenable"
	C_EMU_HOST                = "emuhost"
	C_EMU_HOSTNUM             = "emuhostnum"
	C_EMU_METRIC              = "emumetric"
	C_EMU_METRICNUM           = "emumetricnum"
	C_EMU_TPL                 = "tpl"
	C_EMU_TPLNUM              = "tplnum"
	// falcon-plus
	C_DEBUG                = "debug"
	C_RRD_STORAGE          = "rrd_storage"
	C_MIGRATE_ENABLE       = "migrate_enabled"
	C_MIGRATE_CONCURRENCY  = "migrate_concurrency"
	C_MIGRATE_REPLICAS     = "migrate_replicas"
	C_MIGRATE_CLUSTER      = "migrate_cluster"
	C_MIGRATE_NEW_ENDPOINT = "migrate_newendpoint"
	C_MINSTEP              = "minstep"
	C_SOCKET_ENABLE        = "socket_enable"
	C_SOCKET_ADDR          = "socket_listen"
	C_SOCKET_TIMEOUT       = "socket_timeout"
	C_JUDGE_ENABLE         = "judge_enabled"
	C_JUDGE_BATCH          = "judge_batch"
	C_JUDGE_CONNTIMEOUT    = "judge_conntimeout"
	C_JUDGE_CALLTIMEOUT    = "judge_calltimeout"
	C_JUDGE_MAXCONNS       = "judge_maxconns"
	C_JUDGE_MAXIDLE        = "judge_maxidle"
	C_JUDGE_REPLICAS       = "judge_replicas"
	C_JUDGE_CLUSTER        = "judge_cluster"
	C_GRAPH_ENABLE         = "graph_enabled"
	C_GRAPH_BATCH          = "graph_batch"
	C_GRAPH_CONNTIMEOUT    = "graph_conntimeout"
	C_GRAPH_CALLTIMEOUT    = "graph_calltimeout"
	C_GRAPH_MAXCONNS       = "graph_maxconns"
	C_GRAPH_MAXIDLE        = "graph_maxidle"
	C_GRAPH_REPLICAS       = "graph_replicas"
	C_GRAPH_CLUSTER        = "graph_cluster"
	C_TSDB_ENABLE          = "tsdb_enabled"
	C_TSDB_BATCH           = "tsdb_batch"
	C_TSDB_CONNTIMEOUT     = "tsdb_conntimeout"
	C_TSDB_CALLTIMEOUT     = "tsdb_calltimeout"
	C_TSDB_MAXCONNS        = "tsdb_maxconns"
	C_TSDB_MAXIDLE         = "tsdb_maxidle"
	C_TSDB_RETRY           = "tsdb_retry"
	C_TSDB_ADDRESS         = "tsdb_address"
)

var (
	APP_CONF_NAME = [APP_CONF_SIZE]string{
		"default", "db", "file",
	}

	ConfDefault = map[string]map[string]string{
		"agent": map[string]string{
			C_CONN_TIMEOUT:     "1000",
			C_CALL_TIMEOUT:     "5000",
			C_WORKER_PROCESSES: "2",
			C_HTTP_ENABLE:      "true",
			C_HTTP_ADDR:        "127.0.0.1:1988",
			C_RPC_ENABLE:       "true",
			C_RPC_ADDR:         "127.0.0.1:1989",
			C_GRPC_ENABLE:      "true",
			C_GRPC_ADDR:        "127.0.0.1:1990",
			C_INTERVAL:         "60",
			C_PAYLOADSIZE:      "16",
			C_IFACE_PREFIX:     "eth,em",
		},
		"loadbalance": map[string]string{
			C_CONN_TIMEOUT:     "1000",
			C_CALL_TIMEOUT:     "5000",
			C_WORKER_PROCESSES: "2",
			C_HTTP_ENABLE:      "true",
			C_HTTP_ADDR:        "127.0.0.1:6060",
			C_RPC_ENABLE:       "true",
			C_RPC_ADDR:         "127.0.0.1:8433",
			C_GRPC_ENABLE:      "true",
			C_GRPC_ADDR:        "127.0.0.1:8434",
			C_PAYLOADSIZE:      "16",
		},
		"backend": map[string]string{
			C_CONN_TIMEOUT:     "1000",
			C_CALL_TIMEOUT:     "5000",
			C_WORKER_PROCESSES: "2",
			C_HTTP_ENABLE:      "true",
			C_HTTP_ADDR:        "127.0.0.1:7021",
			C_RPC_ENABLE:       "true",
			C_RPC_ADDR:         "127.0.0.1:7020",
			C_GRPC_ENABLE:      "true",
			C_GRPC_ADDR:        "127.0.0.1:7022",
			C_IDX:              "true",
			C_IDXINTERVAL:      "30",
			C_IDXFULLINTERVAL:  "86400",
			C_DB_MAX_IDLE:      "4",
			C_SHMMAGIC:         "0x80386",
			C_SHMKEY:           "0x7020",
			C_SHMSIZE:          "0x10000000",
		},
		"ctrl": map[string]string{
			C_RUN_MODE:                "dev",
			C_HTTP_ADDR:               "8001",
			C_ENABLE_DOCS:             "true",
			C_SEESION_GC_MAX_LIFETIME: "86400",
			C_SESSION_COOKIE_LIFETIME: "86400",
			C_AUTH_MODULE:             "ldap",
			C_CACHE_MODULE:            "host,role,system,tag,user",
			C_DB_MAX_CONN:             "30",
			C_DB_MAX_IDLE:             "30",
		},
	}
)

type Configer struct {
	data [APP_CONF_SIZE]map[string]string
}

func (c Configer) String() string {
	s := ""
	for i := 0; i < len(c.data); i++ {
		s1 := ""
		for k, v := range c.data[i] {
			s1 += fmt.Sprintf("%-17s %s\n", k, v)
		}
		if len(s1) > 0 {
			s += fmt.Sprintf("%-17s {\n%s\n}\n",
				APP_CONF_NAME[i], IndentLines(1, s1))
		} else {
			s += fmt.Sprintf("%-17s { }\n",
				APP_CONF_NAME[i])
		}
	}
	return s
}

func (c *Configer) Set(model int, m map[string]string) error {
	if model >= APP_CONF_SIZE || model < 0 {
		return errors.New("no model")
	}
	data := make(map[string]string)
	for k, v := range m {
		if len(k) == 0 {
			return errors.New("empty key")
		}
		data[strings.ToLower(k)] = v
	}
	c.data[model] = data
	return nil
}

func (c Configer) Get() [APP_CONF_SIZE]map[string]string {
	return c.data
}

// Bool returns the boolean value for a given key.
func (c *Configer) Bool(key string) (bool, error) {
	return ParseBool(c.getdata(key))
}

// DefaultBool returns the boolean value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultBool(key string, defaultval bool) bool {
	v, err := c.Bool(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int returns the integer value for a given key.
func (c *Configer) Int(key string) (int, error) {
	return strconv.Atoi(c.getdata(key))
}

// DefaultInt returns the integer value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultInt(key string, defaultval int) int {
	v, err := c.Int(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Int64 returns the int64 value for a given key.
func (c *Configer) Int64(key string) (int64, error) {
	return strconv.ParseInt(c.getdata(key), 10, 64)
}

// DefaultInt64 returns the int64 value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultInt64(key string, defaultval int64) int64 {
	v, err := c.Int64(key)
	if err != nil {
		return defaultval
	}
	return v
}

// Float returns the float value for a given key.
func (c *Configer) Float(key string) (float64, error) {
	return strconv.ParseFloat(c.getdata(key), 64)
}

// DefaultFloat returns the float64 value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultFloat(key string, defaultval float64) float64 {
	v, err := c.Float(key)
	if err != nil {
		return defaultval
	}
	return v
}

// String returns the string value for a given key.
func (c *Configer) Str(key string) string {
	return c.getdata(key)
}

// DefaultString returns the string value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultString(key string, defaultval string) string {
	v := c.Str(key)
	if v == "" {
		return defaultval
	}
	return v
}

// Strings returns the []string value for a given key.
// Return nil if config value does not exist or is empty.
func (c *Configer) Strings(key string) []string {
	v := c.Str(key)
	if v == "" {
		return nil
	}
	return strings.Split(v, ";")
}

// DefaultStrings returns the []string value for a given key.
// if err != nil return defaltval
func (c *Configer) DefaultStrings(key string, defaultval []string) []string {
	v := c.Strings(key)
	if v == nil {
		return defaultval
	}
	return v
}

// section.key or key
func (c *Configer) getdata(key string) string {
	if len(key) == 0 {
		return ""
	}

	key = strings.ToLower(key)

	if v, ok := c.data[APP_CONF_FILE][key]; ok {
		return v
	}
	if v, ok := c.data[APP_CONF_DB][key]; ok {
		return v
	}
	if v, ok := c.data[APP_CONF_DEFAULT][key]; ok {
		return v
	}

	return ""
}

// ParseBool returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ParseBool(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%s", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1 {
				return true, nil
			} else if v == 0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

func AssignConfig(ac *Configer, ps ...interface{}) error {
	for _, p := range ps {
		assignSingleConfig(ac, p)
	}
	return nil
}

func assignSingleConfig(ac *Configer, p interface{}) {
	pt := reflect.TypeOf(p)
	if pt.Kind() != reflect.Ptr {
		return
	}
	pt = pt.Elem()
	if pt.Kind() != reflect.Struct {
		return
	}
	pv := reflect.ValueOf(p).Elem()

	for i := 0; i < pt.NumField(); i++ {
		pf := pv.Field(i)
		if !pf.CanSet() {
			continue
		}
		name := pt.Field(i).Name
		switch pf.Kind() {
		case reflect.String:
			pf.SetString(ac.DefaultString(name, pf.String()))
		case reflect.Int, reflect.Int64:
			pf.SetInt(int64(ac.DefaultInt64(name, pf.Int())))
		case reflect.Bool:
			pf.SetBool(ac.DefaultBool(name, pf.Bool()))
		case reflect.Struct:
		default:
			//do nothing here
		}
	}
}
