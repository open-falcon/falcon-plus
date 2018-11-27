// Copyright 2018 Xiaomi, Inc.
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

package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/toolkits/net/httplib"

	coommonModel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
	portal "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
)

type ActionCache struct {
	sync.RWMutex
	M map[int]*portal.Action
}

var Actions = &ActionCache{M: make(map[int]*portal.Action)}

func (this *ActionCache) Get(id int) *portal.Action {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[id]
	if !exists {
		return nil
	}

	return val
}

func (this *ActionCache) Set(id int, action *portal.Action) {
	this.Lock()
	defer this.Unlock()
	this.M[id] = action
}

func getAction(id int) *portal.Action {
	action := CurlAction(id)
	if action != nil {
		Actions.Set(id, action)
	} else {
		action = Actions.Get(id)
	}

	return action
}

func GetAction(event *coommonModel.Event) *portal.Action {
	actionId := event.ActionId()
	action := getAction(actionId)
	if action == nil {
		return nil
	}
	return action
}

func CurlAction(id int) *portal.Action {
	if id <= 0 {
		return nil
	}

	uri := fmt.Sprintf("%s/api/v1/action/%d", config.ApiCon.PlusApi, id)
	req := httplib.Get(uri).SetTimeout(5*time.Second, 30*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm-manager",
		"sig":  config.ApiCon.PlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var act portal.Action
	err := req.ToJson(&act)
	if err != nil {
		log.Errorf("curl %s fail: %v", uri, err)
		return nil
	}

	return &act
}
