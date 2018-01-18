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

package cache

import (
	"log"
	"time"
)

func Init() {
	log.Println("cache begin")

	log.Println("#1 GroupPlugins...")
	GroupPlugins.Init()

	log.Println("#2 GroupTemplates...")
	GroupTemplates.Init()

	log.Println("#3 HostGroupsMap...")
	HostGroupsMap.Init()

	log.Println("#4 HostMap...")
	HostMap.Init()

	log.Println("#5 TemplateCache...")
	TemplateCache.Init()

	log.Println("#6 Strategies...")
	Strategies.Init(TemplateCache.GetMap())

	log.Println("#7 HostTemplateIds...")
	HostTemplateIds.Init()

	log.Println("#8 ExpressionCache...")
	ExpressionCache.Init()

	log.Println("#9 MonitoredHosts...")
	MonitoredHosts.Init()

	log.Println("#10 AgentsInfo...")
	Agents.Init()

	log.Println("cache done")

	go LoopInit()

}

func LoopInit() {
	for {
		time.Sleep(time.Minute)
		GroupPlugins.Init()
		GroupTemplates.Init()
		HostGroupsMap.Init()
		HostMap.Init()
		TemplateCache.Init()
		Strategies.Init(TemplateCache.GetMap())
		HostTemplateIds.Init()
		ExpressionCache.Init()
		MonitoredHosts.Init()
		Agents.Init()
	}
}
