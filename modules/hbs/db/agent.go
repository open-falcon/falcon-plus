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
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/hbs/g"
	"log"
)

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	queryHost := fmt.Sprintf(
		"select count(hostname) as count from host where hostname='%s'",
		agentInfo.ReportRequest.Hostname,
	)
	count := 0
	err := DB.QueryRow(queryHost).Scan(count)
	if err != nil {
		log.Println("query", queryHost, "fail", err)
		return
	}

	sql := ""
	if g.Config().Hosts == "" {
		if count > 0 {
			sql = fmt.Sprintf(
				"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
				agentInfo.ReportRequest.Hostname,
			)
		} else {
			sql = fmt.Sprintf(
				"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s')",
				agentInfo.ReportRequest.Hostname,
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
			)
		}
	} else {
		// sync, just update
		sql = fmt.Sprintf(
			"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.Hostname,
		)
	}

	_, err = DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

}
