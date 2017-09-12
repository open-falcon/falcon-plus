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

func TestAlarmEventNote(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:8088/api/v1/alarm"
	Apitoken := `{"name": "root", "sig": "233fdb00f99811e68a5c001500c6ca5a"}`
	eventId := "s_165_cef145900bf4e2a4a0db8b85762b9cdb"
	rt := resty.New()
	rt.SetHeader("Apitoken", Apitoken)
	Convey("Get notes Test 1, time filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetQueryParams(map[string]string{
				"startTime": "1466611200",
				"endTime":   "1466697600",
			}).
			Get(fmt.Sprintf("%s/event_note", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("event_caseId!='--'")
		// log.Debugf("%v\n", resp.String())
		So(len(check.([]interface{})), ShouldEqual, 2)
	})
	Convey("Get notes Test 1, id filter", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetQueryParams(map[string]string{
				"event_id": eventId,
			}).
			Get(fmt.Sprintf("%s/event_note", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query(fmt.Sprintf("event_caseId='%s'", eventId))
		// log.Debugf("%v\n", resp.String())
		So(len(check.([]interface{})), ShouldBeGreaterThanOrEqualTo, 1)
	})
	Convey("Add Note Test 1, add comment", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`
			{
				"event_id": "%s",
				"note": "test note",
				"status": "comment"
			}`, eventId)).
			Post(fmt.Sprintf("%s/event_note", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query(fmt.Sprintf("id='%s'", eventId))
		So(check, ShouldNotBeNil)
	})
	Convey("Add Note Test 1, change process status", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`
			{
				"event_id": "%s",
				"note": "test note",
				"status": "ignored",
				"case_id": "a000001"
			}`, eventId)).
			Post(fmt.Sprintf("%s/event_note", host))
		log.Debugf("%v", resp.String())
		resp, _ = rt.R().
			SetQueryParams(map[string]string{
				"event_id": eventId,
			}).
			Get(fmt.Sprintf("%s/eventcases", host))
		parser, _ := jsonql.NewStringQuery(resp.String())
		check, _ := parser.Query("status='ignored'")
		So(check, ShouldNotBeNil)
	})

}
