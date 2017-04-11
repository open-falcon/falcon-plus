package redi

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"log"
)

const (
	SMS_QUEUE_NAME  = "/sms"
	MAIL_QUEUE_NAME = "/mail"
)

func PopAllSms() []*model.Sms {
	ret := []*model.Sms{}
	queue := SMS_QUEUE_NAME

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println(err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var sms model.Sms
		err = json.Unmarshal([]byte(reply), &sms)
		if err != nil {
			log.Println(err, reply)
			continue
		}

		ret = append(ret, &sms)
	}

	return ret
}

func PopAllMail() []*model.Mail {
	ret := []*model.Mail{}
	queue := MAIL_QUEUE_NAME

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for {
		reply, err := redis.String(rc.Do("RPOP", queue))
		if err != nil {
			if err != redis.ErrNil {
				log.Println(err)
			}
			break
		}

		if reply == "" || reply == "nil" {
			continue
		}

		var mail model.Mail
		err = json.Unmarshal([]byte(reply), &mail)
		if err != nil {
			log.Println(err, reply)
			continue
		}

		ret = append(ret, &mail)
	}

	return ret
}
