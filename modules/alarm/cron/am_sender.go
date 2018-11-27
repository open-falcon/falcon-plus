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

package cron

import (
	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"net/http"
)

func SendEventToAlarmManager(eve *cmodel.Event) {
	event, err := json.Marshal(eve)
	if err != nil {
		log.Errorf("json marshal err: %s", err.Error())
		return
	}

	url := g.Config().AlarmChannel.AMApi
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(event))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send am fail: %s", err.Error())
		return
	}
	defer resp.Body.Close()
	return
}
