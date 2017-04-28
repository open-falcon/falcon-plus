/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import "fmt"

///////////// CONFIG //////////////
type CmdOpts struct {
	ConfigFile string
	Args       []string
}

type FalconConfig struct {
	ConfigFile string
	pidFile    string
	log        string
	logv       int
	conf       []interface{}
}

func (p FalconConfig) String() string {
	ret := fmt.Sprintf("%-17s %s"+
		"\n%-17s %s"+
		"\n%-17s %d",
		"pidfile", p.pidFile,
		"log", p.log,
		"logv", p.logv,
	)
	for _, v := range p.conf {
		switch GetType(v) {
		case "ConfAgent":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*ConfAgent).Name,
				IndentLines(1, v.(*ConfAgent).String()))
		case "ConfCtrl":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*ConfCtrl).Name,
				IndentLines(1, v.(*ConfCtrl).String()))
		case "ConfLoadbalance":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*ConfLoadbalance).Name,
				IndentLines(1, v.(*ConfLoadbalance).String()))
		case "ConfBackend":
			ret += fmt.Sprintf("\n%s (\n%s\n)",
				v.(*ConfBackend).Name,
				IndentLines(1, v.(*ConfBackend).String()))
		}
	}
	return ret
}

///////////// MODULE //////////////
type ConfAgent struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Configer Configer
}

func (c ConfAgent) String() string {
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		c.Configer.String(),
	)
}

type ConfLoadbalance struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Backend  []LbBackend
	Configer Configer
}

func (c ConfLoadbalance) String() string {
	var s1 string
	for _, v := range c.Backend {
		s1 += fmt.Sprintf("%s\n", v.String())
	}
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"backend", IndentLines(1, s1),
		c.Configer.String(),
	)
}

type LbBackend struct {
	Disabled bool
	Name     string
	Type     string
	Upstream map[string]string
}

func (p LbBackend) String() string {
	var s1, s2 string

	s1 = fmt.Sprintf("%s %s", p.Type, p.Name)
	if p.Disabled {
		s1 += "(Disable)"
	}

	for k, v := range p.Upstream {
		s2 += fmt.Sprintf("%-17s %s\n", k, v)
	}
	return fmt.Sprintf("%s cluster (\n%s\n)", s1, IndentLines(1, s2))
}

type ConfBackend struct {
	Debug    int
	Disabled bool
	Name     string
	Host     string
	Migrate  Migrate
	Configer Configer
}

func (c ConfBackend) String() string {
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"migrate", IndentLines(1, c.Migrate.String()),
		c.Configer.String(),
	)
}

type Migrate struct {
	Disabled bool
	Upstream map[string]string
}

func (p Migrate) String() string {
	var s string

	for k, v := range p.Upstream {
		s += fmt.Sprintf("%-17s %s\n", k, v)
	}
	if s != "" {
		s = fmt.Sprintf("\n%s\n", IndentLines(1, s))
	}

	return fmt.Sprintf("%-17s %v\n"+
		"%s (%s)",
		"disable", p.Disabled,
		"cluster", s)
}

type ConfCtrl struct {
	// only in falcon.conf
	Debug       int
	Disabled    bool
	Name        string
	Host        string
	Metrics     []string
	Ctrl        Configer
	Agent       Configer
	Loadbalance Configer
	Backend     Configer
	Graph       Configer
	Transfer    Configer
	// 1: default, 2: db, 3: ConfCtrl.Container
	// height will cover low
}

func (c ConfCtrl) String() string {
	var s string
	for k, v := range c.Metrics {
		s += fmt.Sprintf("%s ", v)
		if k%5 == 4 {
			s += "\n"
		}
	}
	return fmt.Sprintf("%-17s %d\n"+
		"%-17s %v\n"+
		"%-17s %s\n"+
		"%-17s %s\n"+
		"%s (\n%s\n)\n"+
		"%s",
		"debug", c.Debug,
		"disabled", c.Disabled,
		"Name", c.Name,
		"Host", c.Host,
		"Metrics", IndentLines(1, s),
		c.Ctrl.String(),
	)
}
