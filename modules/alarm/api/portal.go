package api

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/net/httplib"
	"log"
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

	uri := fmt.Sprintf("%s/api/v1/action/%d", g.Config().FalconPlusApi, id)
	req := httplib.Get(uri).SetTimeout(5*time.Second, 30*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  g.Config().FalconPlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var act Action
	err := req.ToJson(&act)
	if err != nil {
		log.Printf("curl %s fail: %v", uri, err)
		return nil
	}

	return &act
}
