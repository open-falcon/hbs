package db

import (
	"fmt"
	"github.com/open-falcon/common/model"
	"github.com/open-falcon/hbs/g"
	"log"
)

func UpdateAgent(agentInfo *model.AgentUpdateInfo) {
	sql := ""
	if g.Config().Hosts == "" {
		sql = fmt.Sprintf(
			"insert into host(hostname, ip, agent_version, plugin_version) values ('%s', '%s', '%s', '%s') on duplicate key update ip='%s', agent_version='%s', plugin_version='%s'",
			agentInfo.ReportRequest.Hostname,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
			agentInfo.ReportRequest.IP,
			agentInfo.ReportRequest.AgentVersion,
			agentInfo.ReportRequest.PluginVersion,
		)
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

	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}

}

func QueryAgentVersion() (map[string]string, error) {
	var ret = make(map[string]string)

	sql := "select hostname,version from agent_version where status = 1"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var hostname, version string
		err = rows.Scan(&hostname, &version)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		ret[hostname] = version
	}

	return ret, nil
}
