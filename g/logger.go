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

import log "github.com/Sirupsen/logrus"
import "github.com/natefinch/lumberjack"
import "os"
import "path"

var (
	logFileName = "logs/" + path.Base(os.Args[0]) + ".log"
	MaxMegaSize = 1000
	MaxBackups  = 5
	MaxAgeDays  = 7
	Compress    = true
)

func InitLog(level string) (err error) {
	rotateLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    MaxMegaSize, // megabytes
		MaxBackups: MaxBackups,  // file numbers
		MaxAge:     MaxAgeDays,  // days
		Compress:   Compress,    // disabled by default
		LocalTime:  true,
	}

	log.SetOutput(rotateLogger)
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.Fatal("log conf only allow [info, debug, warn], please check your confguire")
	}
	return
}
