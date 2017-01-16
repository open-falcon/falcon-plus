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
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:3000/api/v1/graph"
	// cname := "test1"
	// csig := "d4f71cba377911e699d60242ac110010"
	Apitoken := `{"name": "test1", "sig": "d4f71cba377911e699d60242ac110010"}`
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

}
