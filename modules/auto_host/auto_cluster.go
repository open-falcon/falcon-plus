package main

import (
	"fmt"
	"strings"

	//log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/auto_aggr"
	fp "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/graph"
	"log"
)

func genAggregator() {
	endpointCounterList := getEndpointCounters()
	for _, endpointCounter := range endpointCounterList {
		err := genAggr(endpointCounter)
		if err != nil {
			log.Printf("gen cluster for metric(%s) fail:%v", endpointCounter.Counter, err)
			continue
		}
		log.Printf("gen cluster for metric(%s) success", endpointCounter.Counter)
	}

}

func getEndpointCounters() []auto_aggr.EndpointCounter {
	//for get right table name
	enpsHelp := auto_aggr.EndpointCounter{}
	enps := []auto_aggr.EndpointCounter{}
	db.AutoAggr.Table(enpsHelp.TableName()).Scan(&enps)
	return enps
}

func getEndpointName(id uint) (string, error) {
	ep := graph.Endpoint{}
	if err := db.Graph.Table(ep.TableName()).First(&ep, id).Error; err != nil {
		return "", fmt.Errorf("get endpoint name by id(%v) fail :", id, err)
	}
	return ep.Endpoint, nil
}
func genAggr(endpointCounter auto_aggr.EndpointCounter) error {
	ep, err := getEndpointName(uint(endpointCounter.EndpointID))
	if err != nil {
		return err
	}
	grpId, grpEndpointName, err := getGrpinfo(ep)
	if err != nil {
		return err
	}
	orgTags := getOrgTags(endpointCounter.Counter)
	numberator := getNumberator(endpointCounter.Counter)
	denominator := getDenominator(orgTags, endpointCounter.Type)
	metric := getMetric(endpointCounter.Counter)
	tags := getNewTags(endpointCounter.Counter)
	dstype := getDstype(endpointCounter.Type)
	cluster := fp.Cluster{
		GrpId:       grpId,
		Numerator:   numberator,
		Denominator: denominator,
		Endpoint:    grpEndpointName,
		Metric:      metric,
		Tags:        tags,
		DsType:      dstype,
		Step:        endpointCounter.Step,
		Creator:     autoUser,
	}
	if err := addCluster(cluster); err != nil {
		log.Printf("addCluster fail: %s", err)
	}
	return nil
}

func addCluster(c fp.Cluster) error {
	return db.Falcon.Table(c.TableName()).FirstOrCreate(&c, c).Error
}

func getGrpinfo(endpoint string) (int64, string, error) {
	grpName := getLeaderName(endpoint)
	grpId, err := findGrpByLeader(grpName, false)
	if err != nil {
		return -1, "", err
	}
	return grpId, grpName, nil
}

func getNumberator(counter string) string {
	return "$($" + counter + ")"
}

func getDenominator(orgTags, typeStr string) string {
	if strings.Contains(orgTags, "metricType=counter") {
		return "$#"
	}
	return "1"
}

func getMetric(counter string) string {
	list := strings.Split(counter, "/")
	return strings.Join(list[0:len(list)-1], "/")
}

func getOrgTags(counter string) string {
	list := strings.Split(counter, "/")
	if len(list) < 2 {
		return ""
	}
	return list[len(list)-1]
}

func getNewTags(orgTags string) string {
	list := strings.Split(orgTags, ",")
	newList := []string{}
	for _, v := range list {
		if strings.Contains(v, "valueType=count") || strings.Contains(v, "reportType=need_aggr") || strings.Contains(v, "metricType=counter") {
			continue
		}
		newList = append(newList, v)
	}
	return strings.Join(newList, ",")
}

func getDstype(typeStr string) string {
	//return "GAUGE"
	return typeStr
}
