package cron

import (
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

var (
	ChatWorkerChan chan int
	SmsWorkerChan  chan int
	MailWorkerChan chan int
)

func InitSenderWorker() {
	workerConfig := g.Config().Worker
	ChatWorkerChan = make(chan int, workerConfig.Chat)
	SmsWorkerChan = make(chan int, workerConfig.Sms)
	MailWorkerChan = make(chan int, workerConfig.Mail)
}
