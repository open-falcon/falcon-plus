package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/graph"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool = config.Con()

const (
	autoUser = "bot"
)

func getNewHost() []graph.Endpoint {
	//for get right table name
	enpsHelp := graph.Endpoint{}
	enps := []graph.Endpoint{}
	db.AutoAggr.Table(enpsHelp.TableName()).Scan(&enps)

	for _, host := range enps {
		log.Debugf("new endpoint (%+v)", host)
	}
	return enps
}

func deleteFromNewHost(ep graph.Endpoint) error {
	return db.AutoAggr.Table(ep.TableName()).Delete(&ep).Error
}

func InGrp(member, grpId int64) (in bool, err error) {
	newGH := falcon_portal.GrpHost{}
	newGH.HostID = member
	newGH.GrpID = grpId
	if err = db.Falcon.Table(newGH.TableName()).FirstOrInit(&newGH, newGH).Error; err != nil {
		return false, err
	}
	return true, nil
}

func findGrpByLeader(leader string) (grp int64, err error) {
	help := falcon_portal.HostGroup{}
	res := []falcon_portal.HostGroup{}

GETGRP:
	if err = db.Falcon.Table(help.TableName()).Where("grp_name = %s", leader).Scan(&res).Error; err != nil {
		return 0, fmt.Errorf("get host_grp fail:%s", err)
	}
	if len(res) > 0 {
		return res[0].ID, nil
	}

	newHG := falcon_portal.HostGroup{
		Name:       leader,
		CreateUser: autoUser,
	}
	if err = db.Falcon.Table(help.TableName()).Create(newHG).Error; err != nil {
		return -1, fmt.Errorf("create host grp (%s) fail :%s", leader, err)
	}
	id, err := getHostId(leader)
	if err != nil {
		return -1, fmt.Errorf("add leader (%s) to host fail :%s", leader, err)
	}
	if _, err = InGrp(id, newHG.ID); err != nil {
		return -1, fmt.Errorf("insertToGrphost to create new grp_host fail:%s", err)
	}
	goto GETGRP
}

func getHostId(name string) (int64, error) {
	newH := falcon_portal.Host{}
	newH.Hostname = name

	if err := db.Falcon.Table(newH.TableName()).FirstOrInit(&newH, newH).Error; err != nil {
		return -1, err
	}
	return newH.ID, nil
}

func getLeader(host graph.Endpoint) string {
	frag := strings.Split(host.Endpoint, "-")
	return strings.Join(frag[:len(frag)-1], "-")
}

func getHostGrp(host graph.Endpoint) (id int64, err error) {
	leader := getLeader(host)
	return findGrpByLeader(leader)
}

func AutoGenHostGrp() {
	hostList := getNewHost()
	for _, host := range hostList {
		hostGrpId, err := getHostGrp(host)
		if err != nil {
			log.Println(err)
			continue
		}

		hostId, err := getHostId(host.Endpoint)
		if err != nil {
			log.Println(err)
			continue
		}

		_, err = InGrp(hostId, hostGrpId)
		if err != nil {
			log.Println(err)
			continue
		}
		deleteFromNewHost(host)
		log.Println("proccess %s success.", host.Endpoint)
	}
}

func Start() {
	go func() {
		AutoGenHostGrp()
		time.Sleep(time.Minute)
	}()
}
