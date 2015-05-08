package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/open-falcon/sender/g"
	"log"
	"time"
)

var ConnPool *redis.Pool

func InitConnPool() {
	redisConfig := g.Config().Redis

	ConnPool = &redis.Pool{
		MaxIdle:     redisConfig.MaxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisConfig.Addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: PingRedis,
	}
}

func PingRedis(c redis.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}
