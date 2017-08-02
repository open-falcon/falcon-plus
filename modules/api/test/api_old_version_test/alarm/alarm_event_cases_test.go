package test

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
	"github.com/elgs/jsonql"
	"github.com/masato25/resty"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAlarmEventCase(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.InfoLevel)
	host := "http://localhost:8088/api/v1/alarm"
	Apitoken := `{"name": "root", "sig": "233fdb00f99811e68a5c001500c6ca5a"}`
	rt := resty.New()
	rt.SetHeader("Apitoken", Apitoken)
	Convey("Get alarmCase Test 1.1, test status filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	    	"endTime": 1480521600,
	    	"status": "PROBLEM",
	    	"process_status": "ignored,unresolved",
	    	"limit": 10
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='PROBLEM'")
		So(len(check.([]interface{})), ShouldEqual, 2)
		check2, _ := parser.Query("status='OK'")
		So(len(check2.([]interface{})), ShouldEqual, 0)
	})
	Convey("Get alarmCase Test 1.2, test status mutiple filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	    	"endTime": 1480521600,
	    	"status": "PROBLEM,OK",
	    	"process_status": "ignored,unresolved",
	    	"limit": 10
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='PROBLEM'")
		So(len(check.([]interface{})), ShouldEqual, 2)
		check2, _ := parser.Query("status='OK'")
		So(len(check2.([]interface{})), ShouldEqual, 7)
	})
	Convey("Get alarmCase Test 2, test procress_status filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	      "endTime": 1480521600,
	      "status": "PROBLEM",
	      "process_status": "unresolved",
	      "limit": 10
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='PROBLEM'")
		log.Debugf("%v", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 0)
	})
	Convey("Get alarmCase Test 3, test without status filter & limit feature", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	      "endTime": 1480521600,
	      "process_status": "unresolved",
	      "limit": 1
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='OK'")
		log.Debugf("%v", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 1)
	})
	Convey("Get alarmCase Test 4, test timerange filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1477584000,
	      "endTime": 1477670400,
	      "process_status": "unresolved",
	      "limit": 10
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='OK'")
		log.Debugf("%v", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 1)
	})
	Convey("Get alarmCase Test 3, test pagging feature", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	      "endTime": 1480521600,
	      "process_status": "unresolved",
				"page": 0,
	      "limit": 1
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("id='s_322_00c5f5c87a71bd4c686c0a4a0544b719'")
		log.Debugf("%v", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 1)
		resp, _ = rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`
			{
	      "startTime": 1466956800,
	      "endTime": 1480521600,
	      "process_status": "unresolved",
				"page": 1,
	      "limit": 1
	    }`).
			Post(fmt.Sprintf("%s/eventcases", host))
		parser, _ = jsonql.NewStringQuery(resp.String())
		check, _ = parser.Query("id='s_322_00a9b9f9f859d8436f9b643ea5a1fb5e'")
		log.Debugf("%v", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 1)
	})
}
