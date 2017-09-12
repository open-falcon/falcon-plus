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
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
	"time"
)

func ConsumeMail() {
	for {
		L := redi.PopAllMail()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendMailList(L)
	}
}

func SendMailList(L []*model.Mail) {
	for _, mail := range L {
		MailWorkerChan <- 1
		go SendMail(mail)
	}
}

func SendMail(mail *model.Mail) {
	defer func() {
		<-MailWorkerChan
	}()

	url := g.Config().Api.Mail
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Param("tos", mail.Tos)
	r.Param("subject", mail.Subject)
	r.Param("content", mail.Content)
	resp, err := r.String()
	if err != nil {
		log.Errorf("send mail fail, receiver:%s, subject:%s, cotent:%s, error:%v", mail.Tos, mail.Subject, mail.Content, err)
	}

	log.Debugf("send mail:%v, resp:%v, url:%s", mail, resp, url)
}
