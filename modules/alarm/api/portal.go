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

package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/net/httplib"
	"sync"
	"time"
)

//TODO:use api/app/model/falcon_portal/action.go
type Action struct {
	Id                 int    `json:"id"`
	Uic                string `json:"uic"`
	Url                string `json:"url"`
	Callback           int    `json:"callback"`
	BeforeCallbackSms  int    `json:"before_callback_sms"`
	BeforeCallbackMail int    `json:"before_callback_mail"`
	AfterCallbackSms   int    `json:"after_callback_sms"`
	AfterCallbackMail  int    `json:"after_callback_mail"`
}

type ActionCache struct {
	sync.RWMutex
	M map[int]*Action
}

var Actions = &ActionCache{M: make(map[int]*Action)}

func (this *ActionCache) Get(id int) *Action {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[id]
	if !exists {
		return nil
	}

	return val
}

func (this *ActionCache) Set(id int, action *Action) {
	this.Lock()
	defer this.Unlock()
	this.M[id] = action
}

func GetAction(id int) *Action {
	action := CurlAction(id)

	if action != nil {
		Actions.Set(id, action)
	} else {
		action = Actions.Get(id)
	}

	return action
}

func CurlAction(id int) *Action {
	if id <= 0 {
		return nil
	}

	uri := fmt.Sprintf("%s/api/v1/action/%d", g.Config().Api.PlusApi, id)
	req := httplib.Get(uri).SetTimeout(5*time.Second, 30*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  g.Config().Api.PlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var act Action
	err := req.ToJson(&act)
	if err != nil {
		log.Errorf("curl %s fail: %v", uri, err)
		return nil
	}

	return &act
}
