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

package cron

import (
	"encoding/json"
	"time"

	log "github.com/Sirupsen/logrus"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	eventmodel "github.com/open-falcon/falcon-plus/modules/alarm/model/event"
)

func SinglePopEvent(high bool, params ...interface{}) {

	log.Infof("singlePopEvent bool:%t, %v", high, params)

	for {
		reply, err := g.RedisStrings(g.RedisDo("BRPOP", params...))
		if err != nil {
			log.Errorf("brpop alarm event from redis fail: %v, %v", err, params)
			return
		}

		var event cmodel.Event
		err = json.Unmarshal([]byte(reply[1]), &event)
		if err != nil {
			log.Errorf("parse alarojum event fail: %v, %+v", err, reply)
			return
		}

		log.Debugf("pop event: %s", event.String())

		//insert event into database
		eventmodel.InsertEvent(&event)
		// events no longer saved in memory

		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		consume(&event, high)
	}
}
