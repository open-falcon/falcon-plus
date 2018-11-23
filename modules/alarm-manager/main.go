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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/open-falcon/falcon-plus/modules/alarm-manager/config"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/http"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/model/event"
	"github.com/open-falcon/falcon-plus/modules/alarm-manager/model/fault"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configure file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "show help")
	flag.Parse()
	if *version {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	viper.SetConfigFile(*cfg)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config.InitDB(viper.GetViper())
	config.InitApi(viper.GetViper())
	fault.Init()
	event.Init()

	http.Start(viper.GetString("log_level"), viper.GetString("listen"))
}
