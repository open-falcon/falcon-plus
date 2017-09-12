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
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var (
	IMWorkerChan   chan int
	SmsWorkerChan  chan int
	MailWorkerChan chan int
)

func InitSenderWorker() {
	workerConfig := g.Config().Worker
	IMWorkerChan = make(chan int, workerConfig.IM)
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
}
