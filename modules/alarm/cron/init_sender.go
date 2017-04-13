package cron

import (
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var (
	SmsWorkerChan  chan int
	MailWorkerChan chan int
)

func InitSenderWorker() {
	workerConfig := g.Config().Worker
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
}
