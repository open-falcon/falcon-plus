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
	redisConfig := Config().Redis
	if redisConfig.Cluster {
		var err error
		RedisCluster, err = redisCluster.NewCluster(
			&redisCluster.Options{
				StartNodes:   []string{redisConfig.Addr},
				ConnTimeout:  500 * time.Millisecond,
				ReadTimeout:  500 * time.Millisecond,
				WriteTimeout: 500 * time.Millisecond,
				KeepAlive:    16,
				AliveTime:    60 * time.Second,
			})
		if err != nil {
			log.Println("[ERROR] redis cluster init fail", err)
		}
	} else {
		RedisConnPool = &redisgo.Pool{
			MaxIdle:     redisConfig.MaxIdle,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redisgo.Conn, error) {
				c, err := redisgo.Dial("tcp", redisConfig.Addr)
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

func RedisString(reply interface{}, err error) (string, error) {
	redisConfig := Config().Redis
	if redisConfig.Cluster {
		reply, err := redisCluster.String(reply, err)
		if err == redisCluster.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisString", err, reply)
		}
		return reply, err
	} else {
		rc := RedisConnPool.Get()
		defer rc.Close()
		reply, err := redisgo.String(reply, err)
		if err == redisgo.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisString", err, reply)
		}
		return reply, err
	}
}

func RedisStrings(reply interface{}, err error) ([]string, error) {
	redisConfig := Config().Redis
	if redisConfig.Cluster {
		reply, err := redisCluster.Strings(reply, err)
		if err == redisCluster.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisString", err, reply)
		}
		return reply, err
	} else {
		rc := RedisConnPool.Get()
		defer rc.Close()
		reply, err := redisgo.Strings(reply, err)
		if err == redisgo.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisString", err, reply)
		}
		return reply, err
	}
}

func RedisDo(commandName string, args ...interface{}) (interface{}, error) {
	redisConfig := Config().Redis
	if redisConfig.Cluster {
		reply, err := RedisCluster.Do(commandName, args...)
		if err == redisCluster.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisDo", err, commandName, args)
		}
		return reply, err
	} else {
		rc := RedisConnPool.Get()
		defer rc.Close()
		reply, err := rc.Do(commandName, args...)
		if err == redisgo.ErrNil {
			err = nil
		}
		if err != nil {
			log.Println("[ERROR] RedisDo", err, commandName, args)
		}
		return reply, err
	}
}

func RedisClose() {
	redisConfig := Config().Redis
	if redisConfig.Cluster {
		RedisCluster.Close()
	} else {
		RedisConnPool.Close()
	}
}
