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
	"github.com/elgs/gojq"
	"github.com/masato25/resty"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTeam(t *testing.T) {
	viper.AddConfigPath("../../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:3000/api/v1"
	cname := "root"
	csig := "dd81ea033c2d11e6a95d0242ac11000c"
	Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
	teamNmae := "ateamname"
	crateTest := false
	if crateTest {
		Convey("Create Team Failed", t, func() {
			rt := resty.New()
			rt.SetHeader("Apitoken", Apitoken)
			resp, _ := rt.R().
				SetHeader("Content-Type", "application/json").
				SetBody(`{"resume": "i'm descript"}`).
				Post(fmt.Sprintf("%s/team", host))
			So(resp.StatusCode(), ShouldEqual, 400)
		})
		Convey("Create Team Scuessed", t, func() {
			rt := resty.New()
			rt.SetHeader("Apitoken", Apitoken)
			resp, _ := rt.R().
				SetHeader("Content-Type", "application/json").
				SetBody(fmt.Sprintf(`{"team_name": "%s","resume": "i'm descript", "users": [1]}`, teamNmae)).
				Post(fmt.Sprintf("%s/team", host))
			log.Debug(resp.String())
			So(resp.StatusCode(), ShouldEqual, 200)
		})
	} else {
		Convey("Modify Team Group", t, func() {
			Convey("Get Team List", func() {
				rt := resty.New()
				rt.SetHeader("Apitoken", Apitoken)
				resp, _ := rt.R().
					SetQueryParam("q", ".+").
					Get(fmt.Sprintf("%s/team", host))
				So(resp.StatusCode(), ShouldEqual, 200)
			})
			Convey("Get Team List with params", func() {
				var id int
				rt := resty.New()
				rt.SetHeader("Apitoken", Apitoken)
				resp, _ := rt.R().
					SetQueryParam("q", teamNmae).
					Get(fmt.Sprintf("%s/team", host))
				jss, _ := gojq.NewStringQuery(resp.String())
				idtmp, _ := jss.Query("[0].id")
				id = int(idtmp.(float64))
				log.Debug("team id:", id)
				So(resp.StatusCode(), ShouldEqual, 200)
				Convey("Update A Team", func() {
					rt := resty.New()
					rt.SetHeader("Apitoken", Apitoken)
					resp, _ := rt.R().
						SetHeader("Content-Type", "application/json").
						SetBody(fmt.Sprintf(`{"team_id": %d, "resume": "i'm descript update", "users": [4,5,6,7]}`, id)).
						Put(fmt.Sprintf("%s/team", host))
					log.Debugf("reponsed: %v, team_id: %v", resp.String(), id)
					So(resp.StatusCode(), ShouldEqual, 200)
				})
				Convey("Get A Team", func() {
					rt := resty.New()
					rt.SetHeader("Apitoken", Apitoken)
					resp, _ := rt.R().
						Get(fmt.Sprintf("%s/team/%d", host, id))
					log.Debugf("reponsed: %v, team_id: %v", resp.String(), id)
					So(resp.StatusCode(), ShouldEqual, 200)
				})
				Convey("Delete A Team", func() {
					rt := resty.New()
					rt.SetHeader("Apitoken", Apitoken)
					resp, _ := rt.R().
						Delete(fmt.Sprintf("%s/team/%d", host, id))
					log.Debugf("reponsed: %v, team_id: %v", resp.String(), id)
					So(resp.StatusCode(), ShouldEqual, 200)
				})
			})
		})
	}
}
