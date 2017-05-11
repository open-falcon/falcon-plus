package cron

import (
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"github.com/toolkits/net/httplib"
	"time"
)

func ConsumeChat() {
	for {
		L := redi.PopAllChat()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendChatList(L)
	}
}

func SendChatList(L []*model.Chat) {
	for _, chat := range L {
		ChatWorkerChan <- 1
		go SendChat(chat)
	}
}

func SendChat(chat *model.Chat) {
	defer func() {
		<-ChatWorkerChan
	}()

	url := g.Config().Api.Chat
	r := httplib.Post(url).SetTimeout(5*time.Second, 30*time.Second)
	r.Param("tos", chat.Tos)
	r.Param("content", chat.Content)
	resp, err := r.String()
	if err != nil {
		log.Errorf("send chat fail, tos:%s, cotent:%s, error:%v", chat.Tos, chat.Content, err)
	}

	log.Debugf("send chat:%v, resp:%v, url:%s", chat, resp, url)
}
