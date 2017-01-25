package graph

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/sdk/requests"
)

var GraphLastUrl = "http://127.0.0.1:9966/graph/last"

func Last(endpoint, counter string) (val float64, ts int64, err error) {
	param := &model.GraphLastParam{Endpoint: endpoint, Counter: counter}
	bs, err := json.Marshal([]*model.GraphLastParam{param})
	if err != nil {
		return val, ts, err
	}

	bf := bytes.NewBuffer(bs)

	resp, err := http.Post(GraphLastUrl, "application/json", bf)
	if err != nil {
		return val, ts, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return val, ts, err
	}

	var L []*model.GraphLastResp
	err = json.Unmarshal(body, &L)
	if err != nil {
		return val, ts, err
	}

	if len(L) == 0 {
		return val, ts, nil
	}

	v := L[0].Value
	if v == nil {
		return val, ts, nil
	}

	return float64(v.Value), v.Timestamp, nil
}

func Lasts(params []*model.GraphLastParam) ([]*model.GraphLastResp, error) {
	if len(params) == 0 {
		return []*model.GraphLastResp{}, nil
	}

	body, err := requests.PostJsonBody(GraphLastUrl, params)
	if err != nil {
		return []*model.GraphLastResp{}, err
	}

	var L []*model.GraphLastResp
	err = json.Unmarshal(body, &L)
	if err != nil {
		return []*model.GraphLastResp{}, err
	}

	return L, nil
}
