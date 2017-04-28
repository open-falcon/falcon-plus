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
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/golang/glog"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	_ "github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
	_ "github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models/auth"
	_ "github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/routers"
	"github.com/yubo/gotool/flags"
	//_ "github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models/plugin/demo"
)

var opts falcon.CmdOpts

func init() {
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("3")

	flag.StringVar(&opts.ConfigFile, "config",
		"./etc/falcon.conf", "falcon config file")

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s [OPTIONS] COMMAND ",
		"start|stop|reload\n", os.Args[0])
	flags.NewCommand("help", "show help information", help, flag.ExitOnError)
	flags.NewCommand("start", "start falcon", start, flag.ExitOnError)
	flags.NewCommand("stop", "stop falcon", stop, flag.ExitOnError)
	flags.NewCommand("parse", "just parse falcon ConfigFile", parse, flag.ExitOnError)
	flags.NewCommand("reload", "reload falcon", reload, flag.ExitOnError)
}

func help(arg interface{}) {
	flags.Usage()
}

func start(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err := app.Save(); err != nil {
		glog.Fatal(err)
	}

	dir, _ := os.Getwd()
	glog.V(4).Infof("work dir :%s", dir)
	glog.V(4).Infof("\n%s", c)

	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Start()
}

func stop(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func parse(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	conf := falcon.Parse(opts.ConfigFile, true)
	dir, _ := os.Getwd()
	glog.Infof("work dir :%s", dir)
	glog.Infof("\n%s", conf)
}

func reload(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}

func main() {
	flags.Parse()
	cmd := flags.CommandLine.Cmd

	if cmd != nil && cmd.Action != nil {
		opts.Args = cmd.Flag.Args()
		cmd.Action(&opts)
	} else {
		//flags.Usage()
		opts.Args = flag.Args()
		start(&opts)
	}
}
