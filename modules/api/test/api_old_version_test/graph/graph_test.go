package test

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
	"github.com/masato25/resty"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraph(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg.example")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:8080/api/v1/graph"
	Apitoken := `{"name": "laiwei4", "sig": "022a294d108e11e7b7d1f45c89cb3693"}`
	rt := resty.New()
	rt.SetHeader("Apitoken", Apitoken)
	Convey("Get Endpoint Failed", t, func() {
		resp, _ := rt.R().Get(fmt.Sprintf("%s/endpoint", host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("Get Endpoint without login session", t, func() {
		resp, _ := resty.R().Get(fmt.Sprintf("%s/endpoint", host))
		So(resp.StatusCode(), ShouldEqual, 401)
	})

	Convey("Get Endpoint List", t, func() {
		resp, _ := rt.R().SetQueryParam("q", "a.+").Get(fmt.Sprintf("%s/endpoint", host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Get Counter Failed", t, func() {
		resp, _ := rt.R().Get(fmt.Sprintf("%s/endpoint_counter", host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Get Counter List", t, func() {
		resp, _ := rt.R().SetQueryParam("eid", "6,7").SetQueryParam("metricQuery", "disk.+").Get(fmt.Sprintf("%s/endpoint_counter", host))
		log.Debug(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Delete counter", t, func() {
		resp, _ := rt.R().SetHeader("Content-Type", "application/json").
			SetBody(`{"endpoints":["laiwei-aggregator-1"], "counters":["agent.alive.percent/name=xx"]}`).
			Delete(fmt.Sprintf("%s/counter", host))
		log.Debug(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Delete endpoint", t, func() {
		resp, _ := rt.R().SetHeader("Content-Type", "application/json").
			SetBody(`["0.0.0.0"]`).
			Delete(fmt.Sprintf("%s/endpoint", host))
		log.Debug(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})

}
