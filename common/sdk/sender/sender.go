package sender

import (
	"log"
	"time"
)

const LIMIT = 200

var MetaDataQueue = NewSafeLinkedList()
var PostPushUrl string
var Debug bool

func StartSender() {
	go startSender()
}

func startSender() {
	for {
		L := MetaDataQueue.PopBack(LIMIT)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}

		err := PostPush(L)
		if err != nil {
			log.Println("[E] push to transfer fail", err)
		}
	}
}
