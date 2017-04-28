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
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
)

var (
	falconModules map[string]Module
)

type Module interface {
	New(config interface{}) Module
	Prestart() error
	Start() error
	Stop() error
	Reload(config interface{}) error
	Signal(os.Signal) error
	String() string
	Name() string
}

// reload not support add/del/disable module
type Process struct {
	Config *FalconConfig
	Pid    int
	status uint32
	module []Module
}

func NewProcess(c *FalconConfig) *Process {
	p := &Process{
		Config: c,
		Pid:    os.Getpid(),
		status: APP_STATUS_PENDING,
		module: make([]Module, len(c.conf)),
	}
	return p
}

func (p *Process) Status() uint32 {
	return atomic.LoadUint32(&p.status)
}

func (p *Process) Kill(sig syscall.Signal) error {
	if pid, err := readFileInt(p.Config.pidFile); err != nil {
		return err
	} else {
		glog.Infof("kill %d %s\n", pid, sig)
		return syscall.Kill(pid, sig)
	}
}

func (p *Process) Check() error {
	pid, err := readFileInt(p.Config.pidFile)
	if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}

	_, err = os.Stat(fmt.Sprintf("/proc/%s", pid))
	if os.IsNotExist(err) {
		return nil
	}

	return fmt.Errorf("proccess %s exist", pid)
}

func (p *Process) Save() error {
	return ioutil.WriteFile(p.Config.pidFile,
		[]byte(fmt.Sprintf("%d", p.Pid)), 0644)
}

func (p *Process) Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1)
	atomic.StoreUint32(&p.status, APP_STATUS_RUNNING)

	setGlog(p.Config)

	for i := 0; i < len(p.Config.conf); i++ {
		m, ok := falconModules[GetType(p.Config.conf[i])]
		if !ok {
			glog.Exitf("%s's module not support", GetType(p.Config.conf[i]))
		}
		p.module[i] = m.New(p.Config.conf[i])
	}

	for i := 0; i < len(p.module); i++ {
		p.module[i].Prestart()
	}

	for i := 0; i < len(p.module); i++ {
		p.module[i].Start()
	}

	glog.Infof(MODULE_NAME+"[%d] register signal notify", p.Pid)

	for {
		s := <-sigs
		glog.Infof(MODULE_NAME+"recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			pidfile := fmt.Sprintf("%s.%d", p.Config.pidFile, p.Pid)
			glog.Info(MODULE_NAME + "exiting")
			atomic.StoreUint32(&p.status, APP_STATUS_EXIT)
			os.Rename(p.Config.pidFile, pidfile)

			for i, n := 0, len(p.module); i < n; i++ {
				p.module[n-i-1].Stop()
			}

			glog.Infof(MODULE_NAME+"pid:%d exit", p.Pid)
			os.Remove(pidfile)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Info(MODULE_NAME + "reload")

			// reparse config, get new config
			newConfig := Parse(p.Config.ConfigFile, false)

			// check config diff
			if len(newConfig.conf) != len(p.Config.conf) {
				glog.Error("not support add/del module\n")
				break
			}

			for i, config := range newConfig.conf {
				m, ok := falconModules[GetType(config)]
				if !ok {
					glog.Exitf("%s's module not support", GetType(config))
					break
				}
				newM := m.New(config)
				if newM.Name() != p.module[i].Name() {
					glog.Exitf("%s's module not support, not support "+
						"add/del/disable module", GetType(config))
					break
				}
			}

			// do it
			atomic.StoreUint32(&p.status, APP_STATUS_RELOAD)
			setGlog(newConfig)
			for i, m := range p.module {
				m.Reload(newConfig.conf[i])
			}
			atomic.StoreUint32(&p.status, APP_STATUS_RUNNING)
		default:
			for _, m := range p.module {
				m.Signal(s)
			}
		}
	}
}

func RegisterModule(name string, m Module) error {
	if _, ok := falconModules[name]; ok {
		return ErrExist
	} else {
		falconModules[name] = m
		return nil
	}
}

func setGlog(c *FalconConfig) {
	glog.V(3).Infof("set glog %s, %d", c.log, c.logv)
	flag.Lookup("v").Value.Set(fmt.Sprintf("%d", c.logv))

	if strings.ToLower(c.log) == "stdout" {
		flag.Lookup("logtostderr").Value.Set("true")
		return
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
	}

	if fi, err := os.Stat(c.log); err != nil || !fi.IsDir() {
		glog.Errorf("log dir %s does not exist or not dir", c.log)
	} else {
		flag.Lookup("logtostderr").Value.Set("false")
		flag.Lookup("log_dir").Value.Set(c.log)
	}
}
