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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/falcon-plus/modules/alarm/api"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"strings"
	"time"
)

func CombineSms() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineSms()
	}
}

func CombineMail() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineMail()
	}
}

func CombineIM() {
	for {
		// 每分钟读取处理一次
		time.Sleep(time.Minute)
		combineIM()
	}
}

func combineMail() {
	dtos := popAllMailDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*MailDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Email, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*MailDto{dtos[i]}
		}
	}

	// 不要在这处理，继续写回redis，否则重启alarm很容易丢数据
	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteMail([]string{arr[0].Email}, arr[0].Subject, arr[0].Content)
			continue
		}

		subject := fmt.Sprintf("[P%d][%s] %d %s", arr[0].Priority, arr[0].Status, size, arr[0].Metric)
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, "\r\n")

		log.Debugf("combined mail subject:%s, content:%s", subject, content)
		redi.WriteMail([]string{arr[0].Email}, subject, content)
	}
}

func combineIM() {
	dtos := popAllImDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*ImDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].IM, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*ImDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteIM([]string{arr[0].IM}, arr[0].Content)
			continue
		}

		// 把多个im内容写入数据库，只给用户提供一个链接
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, ",,")

		first := arr[0].Content
		t := strings.Split(first, "][")
		eg := ""
		if len(t) >= 3 {
			eg = t[2]
		}

		path, err := api.LinkToSMS(content)
		chat := ""
		if err != nil || path == "" {
			chat = fmt.Sprintf("[P%d][%s] %d %s.  e.g. %s. detail in email", arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg)
			log.Error("create short link fail", err)
		} else {
			chat = fmt.Sprintf("[P%d][%s] %d %s e.g. %s %s/portal/links/%s ",
				arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg, g.Config().Api.Dashboard, path)
			log.Debugf("combined im is:%s", chat)
		}

		redi.WriteIM([]string{arr[0].IM}, chat)
	}
}

func combineSms() {
	dtos := popAllSmsDto()
	count := len(dtos)
	if count == 0 {
		return
	}

	dtoMap := make(map[string][]*SmsDto)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d%s%s%s", dtos[i].Priority, dtos[i].Status, dtos[i].Phone, dtos[i].Metric)
		if _, ok := dtoMap[key]; ok {
			dtoMap[key] = append(dtoMap[key], dtos[i])
		} else {
			dtoMap[key] = []*SmsDto{dtos[i]}
		}
	}

	for _, arr := range dtoMap {
		size := len(arr)
		if size == 1 {
			redi.WriteSms([]string{arr[0].Phone}, arr[0].Content)
			continue
		}

		// 把多个sms内容写入数据库，只给用户提供一个链接
		contentArr := make([]string, size)
		for i := 0; i < size; i++ {
			contentArr[i] = arr[i].Content
		}
		content := strings.Join(contentArr, ",,")

		first := arr[0].Content
		t := strings.Split(first, "][")
		eg := ""
		if len(t) >= 3 {
			eg = t[2]
		}

		path, err := api.LinkToSMS(content)
		sms := ""
		if err != nil || path == "" {
			sms = fmt.Sprintf("[P%d][%s] %d %s.  e.g. %s. detail in email", arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg)
			log.Error("get short link fail", err)
		} else {
			sms = fmt.Sprintf("[P%d][%s] %d %s e.g. %s %s/portal/links/%s ",
				arr[0].Priority, arr[0].Status, size, arr[0].Metric, eg, g.Config().Api.Dashboard, path)
			log.Debugf("combined sms is:%s", sms)
		}

		redi.WriteSms([]string{arr[0].Phone}, sms)
	}
}

func popAllSmsDto() []*SmsDto {
	ret := []*SmsDto{}
	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get SmsDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var smsDto SmsDto
		err = json.Unmarshal([]byte(reply), &smsDto)
		if err != nil {
			log.Errorf("json unmarshal SmsDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &smsDto)
	}

	return ret
}

func popAllMailDto() []*MailDto {
	ret := []*MailDto{}
	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get MailDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mailDto MailDto
		err = json.Unmarshal([]byte(reply), &mailDto)
		if err != nil {
			log.Errorf("json unmarshal MailDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &mailDto)
	}

	return ret
}

func popAllImDto() []*ImDto {
	ret := []*ImDto{}
	queue := g.Config().Redis.UserIMQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Error("get ImDto fail", err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var imDto ImDto
		err = json.Unmarshal([]byte(reply), &imDto)
		if err != nil {
			log.Errorf("json unmarshal imDto: %s fail: %v", reply, err)
			continue
		}

		ret = append(ret, &imDto)
	}

	return ret
}
