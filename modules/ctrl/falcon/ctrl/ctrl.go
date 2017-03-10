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
package ctrl

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
)

const (
	MODULE_NAME     = "\x1B[32m[CTRL]\x1B[0m "
	CONN_RETRY      = 2
	DEBUG_STAT_STEP = 60
)

var (
	prestartHooks = make([]hookfunc, 0)
	reloadHooks   = make([]hookfunc, 0)
	Configure     *falcon.ConfCtrl
	EtcdCli       *falcon.EtcdCli
)

type hookfunc func(conf *falcon.ConfCtrl) error

type Ctrl struct {
	Conf *falcon.ConfCtrl
	// runtime
	status       uint32
	running      chan struct{}
	rpcListener  *net.TCPListener
	httpListener *net.TCPListener
	httpMux      *http.ServeMux
}

func init() {
	falcon.RegisterModule(falcon.GetType(falcon.ConfCtrl{}), &Ctrl{})
}

func RegisterPrestart(fn hookfunc) {
	prestartHooks = append(prestartHooks, fn)
}

func RegisterReload(fn hookfunc) {
	reloadHooks = append(reloadHooks, fn)
}

func (p *Ctrl) New(conf interface{}) falcon.Module {
	return &Ctrl{Conf: conf.(*falcon.ConfCtrl)}
}

func (p *Ctrl) Name() string {
	return fmt.Sprintf("%s", p.Conf.Name)
}

func (p *Ctrl) String() string {
	return p.Conf.String()
}

// ugly hack
// should called by main package
func (p *Ctrl) Prestart() error {
	Configure = p.Conf

	EtcdCli = falcon.NewEtcdCli(Configure.Ctrl)

	EtcdCli.Prestart()
	for _, fn := range prestartHooks {
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}
	return nil
}

func (p *Ctrl) Start() error {
	glog.V(3).Infof(MODULE_NAME+"%s Start()", p.Conf.Name)

	p.status = falcon.APP_STATUS_PENDING
	p.running = make(chan struct{}, 0)

	EtcdCli.Start()
	// p.rpcStart()
	// p.httpStart()
	p.statStart()
	go beego.Run()
	return nil
}

func (p *Ctrl) Stop() error {
	glog.V(3).Infof(MODULE_NAME+"%s Stop()", p.Conf.Name)
	p.status = falcon.APP_STATUS_EXIT
	close(p.running)
	p.statStop()
	EtcdCli.Stop()
	// p.httpStop()
	// p.rpcStop()
	return nil
}

// TODO: reload is not yet implemented
func (p Ctrl) Reload(config interface{}) error {
	p.Conf = config.(*falcon.ConfCtrl)
	glog.V(3).Infof(MODULE_NAME+"%s Reload()", p.Conf.Name)

	Configure = p.Conf

	EtcdCli.Reload(Configure.Ctrl)
	for _, fn := range prestartHooks {
		if err := fn(Configure); err != nil {
			panic(err)
		}
	}

	return nil
}

func (p Ctrl) Signal(sig os.Signal) error {
	glog.V(3).Infof(MODULE_NAME+"%s signal %v", p.Conf.Name, sig)
	return nil
}
