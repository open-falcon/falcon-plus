package sdk

import (
	"encoding/json"
	"fmt"
	cmodel "github.com/open-falcon/common/model"
	"github.com/open-falcon/falcon-plus/common/sdk/requests"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	"github.com/toolkits/net/httplib"
)

func HostnamesByID(group_id int64) ([]string, error) {

	uri := fmt.Sprintf("%s/api/v1/hostgroup/%d", g.Config().Api.PlusApi, group_id)
	req, err := requests.CurlPlus(uri, "GET", "aggregator", g.Config().Api.PlusApiToken,
		map[string]string{}, map[string]string{})

	if err != nil {
		return []string{}, err
	}

	type RESP struct {
		HostGroup f.HostGroup `json:"hostgroup"`
		Hosts     []f.Host    `json:"hosts"`
	}

	resp := &RESP{}
	err = req.ToJson(&resp)
	if err != nil {
		return []string{}, err
	}

	hosts := []string{}
	for _, x := range resp.Hosts {
		hosts = append(hosts, x.Hostname)
	}
	return hosts, nil
}

func QueryLastPoints(endpoints, counters []string) (resp []*cmodel.GraphLastResp, err error) {
	uri := fmt.Sprintf("%s/api/v1/graph/lastpoint", g.Config().Api.PlusApi)

	var req *httplib.BeegoHttpRequest
	headers := map[string]string{"Content-type": "application/json"}
	req, err = requests.CurlPlus(uri, "POST", "aggregator", g.Config().Api.PlusApiToken,
		headers, map[string]string{})

	if err != nil {
		return
	}

	body := map[string][]string{
		"endpoints": endpoints,
		"counters":  counters,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return
	}

	req.Body(b)

	err = req.ToJson(&resp)
	if err != nil {
		return
	}

	return resp, nil
}
