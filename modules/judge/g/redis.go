package g

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var RedisConnPool *redis.Pool

func InitRedisConnPool() {
	if !Config().Alarm.Enabled {
		return
	}

	dsn := Config().Alarm.Redis.Dsn
	maxIdle := Config().Alarm.Redis.MaxIdle
	idleTimeout := 240 * time.Second

	connTimeout := time.Duration(Config().Alarm.Redis.ConnTimeout) * time.Millisecond
	readTimeout := time.Duration(Config().Alarm.Redis.ReadTimeout) * time.Millisecond
	writeTimeout := time.Duration(Config().Alarm.Redis.WriteTimeout) * time.Millisecond

	RedisConnPool = &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", dsn, connTimeout, readTimeout, writeTimeout)
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
