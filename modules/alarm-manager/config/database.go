// Copyright 2018 Xiaomi, Inc.
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

package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type DataBasePool struct {
	AM *gorm.DB
}

var (
	Dbpool DataBasePool
)

func Con() DataBasePool {
	return Dbpool
}

func InitDB(vip *viper.Viper) (err error) {
	log.Println("InitDB begin")
	var p *sql.DB
	am, err := gorm.Open("mysql", vip.GetString("db.alarm_manager"))
	am.Dialect().SetDB(p)
	am.LogMode(viper.GetBool("db.debug"))
	if err != nil {
		return fmt.Errorf("connect to alarmmanager db error: %s", err.Error())
	}
	Dbpool.AM = am
	return
}

func CloseDB() (err error) {
	err = Dbpool.AM.Close()
	if err != nil {
		return
	}
	return
}
