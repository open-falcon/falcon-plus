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

func QueryAgentsInfo() (map[string]*model.AgentUpdateInfo, error) {
	m := make(map[string]*model.AgentUpdateInfo)

	sql := "select hostname, ip, agent_version, plugin_version, update_at from host"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return m, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			hostname       string
			ip             string
			agent_version  string
			plugin_version string
			update_at      int64
		)

		err = rows.Scan(&hostname, &ip, &agent_version, &plugin_version, &update_at)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		m[hostname] = &model.AgentUpdateInfo{
			LastUpdate: update_at,
			ReportRequest: &model.AgentReportRequest{
				Hostname:      hostname,
				IP:            ip,
				AgentVersion:  agent_version,
				PluginVersion: plugin_version,
			},
		}
	}

	return m, nil
}

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	var (
		hostname       string
		ip             string
		agent_version  string
		plugin_version string
	)

	sql := fmt.Sprintf(
		"select hostname, ip, agent_version, plugin_version from host where hostname = %s",
		agentInfo.ReportRequest.Hostname,
	)

	rows, err := DB.Query(sql)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&hostname, &ip, &agent_version, &plugin_version)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
	}

	if agentInfo.ReportRequest.Hostname == hostname && agentInfo.ReportRequest.IP == ip && agentInfo.ReportRequest.AgentVersion == agent_version && agentInfo.ReportRequest.PluginVersion == plugin_version {
		return
	}

	sql = ""
	if g.Config().Hosts == "" {
		if hostname == "" && ip == "" && agent_version == "" && plugin_version == "" {
			sql = fmt.Sprintf(
				"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s')",
				agentInfo.ReportRequest.Hostname,
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
			)
		} else {
			sql = fmt.Sprintf(
				"update host set ip='%s', agent_version='%s', plugin_version='%s' where hostname='%s'",
				agentInfo.ReportRequest.IP,
				agentInfo.ReportRequest.AgentVersion,
				agentInfo.ReportRequest.PluginVersion,
				agentInfo.ReportRequest.Hostname,
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
