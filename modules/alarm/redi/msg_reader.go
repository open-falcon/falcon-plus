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

package redi

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
)

const (
	IM_QUEUE_NAME   = "/im"
	SMS_QUEUE_NAME  = "/sms"
	MAIL_QUEUE_NAME = "/mail"
)

func PopAllSms() []*model.Sms {
	ret := []*model.Sms{}
	queue := SMS_QUEUE_NAME

	for {
		reply, err := g.RedisString(g.RedisDo("RPOP", queue))
		if err != nil {
			log.Error("rpop all sms msg fail:", err)
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var sms model.Sms
		err = json.Unmarshal([]byte(reply), &sms)
		if err != nil {
			log.Error(err, reply)
			continue
		}

		ret = append(ret, &sms)
	}

	return ret
}

func PopAllIM() []*model.IM {
	ret := []*model.IM{}
	queue := IM_QUEUE_NAME

	for {
		reply, err := g.RedisString(g.RedisDo("RPOP", queue))
		if err != nil {
			log.Error("rpop all im msg fail:", err)
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var im model.IM
		err = json.Unmarshal([]byte(reply), &im)
		if err != nil {
			log.Error(err, reply)
			continue
		}

		ret = append(ret, &im)
	}

	return ret
}

func PopAllMail() []*model.Mail {
	ret := []*model.Mail{}
	queue := MAIL_QUEUE_NAME

	for {
		reply, err := g.RedisString(g.RedisDo("RPOP", queue))
		if err != nil {
			log.Error("rpop all mail msg fail:", err)
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mail model.Mail
		err = json.Unmarshal([]byte(reply), &mail)
		if err != nil {
			log.Error(err, reply)
			continue
		}

		ret = append(ret, &mail)
	}

	return ret
}
