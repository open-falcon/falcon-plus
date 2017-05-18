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

package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"log"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", g.Config().Database.Addr)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	DB.SetMaxIdleConns(g.Config().Database.Idle)

	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
