package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/sender/model"
	"log"
)

func PopAllSms(queue string) []*model.Sms {
	ret := []*model.Sms{}

	rc := ConnPool.Get()
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

func PopAllMail(queue string) []*model.Mail {
	ret := []*model.Mail{}

	rc := ConnPool.Get()
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
