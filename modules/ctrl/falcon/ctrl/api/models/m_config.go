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

	"github.com/astaxie/beego/orm"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl"
)

type Kv struct {
	Key     string
	Section string
	Value   string
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

func (op *Operator) SetDbConfig(module string, conf map[string]string) error {
	kv := make(map[string]string)
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
	case "ctrl", "agent", "lb", "backend":
		return op.SetDbConfig(module, conf)
	default:
		return ErrNoModule
	}
}
