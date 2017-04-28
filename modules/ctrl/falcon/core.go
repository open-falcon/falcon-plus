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
package falcon

import (
	"flag"
	"fmt"

	"github.com/yubo/gotool/flags"
)

const (
	IndentSize   = 4
	DEFAULT_STEP = 60 //s
	MIN_STEP     = 30 //s
	VERSION      = "0.0.2"
	REPLICAS     = 500
	GAUGE        = "GAUGE"
	DERIVE       = "DERIVE"
	COUNTER      = "COUNTER"
	MODULE_NAME  = "\x1B[32m[FALCON]\x1B[0m "
)

const (
	APP_STATUS_INIT = iota
	APP_STATUS_PENDING
	APP_STATUS_RUNNING
	APP_STATUS_EXIT
	APP_STATUS_RELOAD
)

func init() {
	falconModules = make(map[string]Module)

	flags.NewCommand("version", "show falcon version information",
		Version, flag.ExitOnError)

	flags.NewCommand("git", "show falcon git version information",
		Git, flag.ExitOnError)

	flags.NewCommand("changelog", "show falcon changelog information",
		Changelog, flag.ExitOnError)
}

func Version(arg interface{}) {
	fmt.Println(VERSION)
}

func Git(arg interface{}) {
	fmt.Println(COMMIT)
}

func Changelog(arg interface{}) {
	fmt.Println(CHANGELOG)
}
