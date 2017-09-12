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

func TestAlarmEvents(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:8088/api/v1/alarm"
	Apitoken := `{"name": "root", "sig": "233fdb00f99811e68a5c001500c6ca5a"}`
	eventId := "s_165_cef145900bf4e2a4a0db8b85762b9cdb"
	eventId2 := "s_322_00f5a68989281d30fbbb647194a09230"
	Convey("Get alarms Test 1, time filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetFormData(map[string]string{
				"startTime": "1466611200",
				"endTime":   "1466628960",
				"event_id":  eventId,
			}).
			Post(fmt.Sprintf("%s/events", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("id!=-1")
		// log.Debugf("%v\n", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 2)
	})
	Convey("Get alarms Test 2, test status filter & limit", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetFormData(map[string]string{
				"event_id": eventId2,
				"status":   "1",
				"limit":    "2",
			}).
			Post(fmt.Sprintf("%s/events", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status=1")
		// log.Debugf("%v\n", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 2)
	})
	Convey("Get alarms Test 3, test pagging", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetFormData(map[string]string{
				"event_id": eventId2,
				"limit":    "1",
				"page":     "0",
			}).
			Post(fmt.Sprintf("%s/events", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("id=2283581")
		So(len(check.([]interface{})), ShouldEqual, 1)
		resp2, _ := rt.R().
			SetFormData(map[string]string{
				"event_id": eventId2,
				"limit":    "1",
				"page":     "1",
			}).
			Post(fmt.Sprintf("%s/events", host))
		parser, _ = jsonql.NewStringQuery(resp2.String())
		check, _ = parser.Query("id=2283566")
		So(len(check.([]interface{})), ShouldEqual, 1)
	})
}
