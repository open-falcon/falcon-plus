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

package g

import (
	redisCluster "github.com/chasex/redis-go-cluster"
	redisgo "github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var RedisConnPool *redisgo.Pool
var RedisCluster *redisCluster.Cluster

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

	if Config().Alarm.Redis.Cluster {
		var err error
		RedisCluster, err = redisCluster.NewCluster(
			&redisCluster.Options{
				StartNodes:   []string{dsn},
				ConnTimeout:  connTimeout,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
				KeepAlive:    16,
				AliveTime:    60 * time.Second,
			})
		if err != nil {
			log.Println("[ERROR] redis cluster init fail", err)
		}
	} else {
		RedisConnPool = &redisgo.Pool{
			MaxIdle:     maxIdle,
			IdleTimeout: idleTimeout,
			Dial: func() (redisgo.Conn, error) {
				c, err := redisgo.DialTimeout("tcp", dsn, connTimeout, readTimeout, writeTimeout)
				if err != nil {
					log.Println("[ERROR] redis pool init fail", err)
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: PingRedis,
		}
	}

}

func PingRedis(c redisgo.Conn, t time.Time) error {
	_, err := c.Do("ping")
	if err != nil {
		log.Println("[ERROR] ping redis fail", err)
	}
	return err
}

func RedisDo(commandName string, args ...interface{}) (interface{}, error) {
	if Config().Alarm.Redis.Cluster {
		reply, err := RedisCluster.Do(commandName, args...)
		if err == redisCluster.ErrNil {
			err = nil
		}
		return reply, err
	} else {
		rc := RedisConnPool.Get()
		defer rc.Close()
		reply, err := rc.Do(commandName, args...)
		if err == redisgo.ErrNil {
			err = nil
		}
		return reply, err
	}
}

func RedisClose() {
	if Config().Alarm.Redis.Cluster {
		RedisCluster.Close()
	} else {
		RedisConnPool.Close()
	}
}
