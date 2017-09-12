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
	"strings"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
)

func lpush(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Error("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteSmsModel(sms *model.Sms) {
	if sms == nil {
		return
	}

	bs, err := json.Marshal(sms)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write sms to queue, sms:%v, queue:%s", sms, SMS_QUEUE_NAME)
	lpush(SMS_QUEUE_NAME, string(bs))
}

func WriteIMModel(im *model.IM) {
	if im == nil {
		return
	}

	bs, err := json.Marshal(im)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write im to queue, im:%v, queue:%s", im, IM_QUEUE_NAME)
	lpush(IM_QUEUE_NAME, string(bs))
}

func WriteMailModel(mail *model.Mail) {
	if mail == nil {
		return
	}

	bs, err := json.Marshal(mail)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("write mail to queue, mail:%v, queue:%s", mail, MAIL_QUEUE_NAME)
	lpush(MAIL_QUEUE_NAME, string(bs))
}

func WriteSms(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	sms := &model.Sms{Tos: strings.Join(tos, ","), Content: content}
	WriteSmsModel(sms)
}

func WriteIM(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	im := &model.IM{Tos: strings.Join(tos, ","), Content: content}
	WriteIMModel(im)
}

func WriteMail(tos []string, subject, content string) {
	if len(tos) == 0 {
		return
	}

	mail := &model.Mail{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	WriteMailModel(mail)
}
