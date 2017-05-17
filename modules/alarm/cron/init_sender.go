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
