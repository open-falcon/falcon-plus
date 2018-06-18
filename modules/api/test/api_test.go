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
	"encoding/json"
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/masato25/resty"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	cfg "github.com/open-falcon/falcon-plus/modules/api/config"
)

var (
	api_host           = "http://localhost:8080/api/v1"
	test_user_name     = "apitest-user1"
	test_user_password = "password"
)

func init() {
	viper.AddConfigPath("../")
	viper.AddConfigPath("./")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	api_host = fmt.Sprintf("http://localhost%s/api/v1", viper.GetString("web_port"))

	if err := cfg.InitDB(viper.GetBool("db.db_bug"), viper.GetViper()); err != nil {
		log.Fatal(err.Error())
	}

	init_testing_user()
}

func init_testing_user() {
	password := utils.HashIt(test_user_password)
	user := uic.User{
		Name:   test_user_name,
		Passwd: password,
		Cnname: test_user_name,
		Email:  test_user_name + "@test.com",
		Phone:  "1234567890",
		IM:     "hellotest",
		QQ:     "380511212",
	}

	db := cfg.Con()
	if db.Uic.Table("user").Where("name = ?", test_user_name).First(&uic.User{}).RecordNotFound() {
		if err := db.Uic.Table("user").Create(&user).Error; err != nil {
			log.Fatal(err)
		}
		log.Info("create_user:", test_user_name)
	}
}

func get_session_token() (string, error) {
	rt := resty.New()
	resp, _ := rt.R().SetQueryParam("name", test_user_name).
		SetQueryParam("password", test_user_password).
		Post(fmt.Sprintf("%s/user/login", api_host))

	type Resp struct {
		Sig   string `json:"sig"`
		Name  string `json:"name"`
		Admin bool   `json:"admin"`
	}
	resp_obj := Resp{}
	err := json.Unmarshal([]byte(resp.String()), &resp_obj)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, test_user_name, resp_obj.Sig)

	return Apitoken, err
}

func TestNodata(t *testing.T) {
	apitoken, _ := get_session_token()
	rt := resty.New()
	rt.SetHeader("Apitoken", apitoken)

	var nid int = 0

	Convey("Create nodata config", t, func() {
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"tags": "", "step": 60, "obj_type": "host", "obj": "docker-agent", "name": "api.testnodata", "mock": -1, "metric": "api.test.metric", "dstype": "GAUGE"}`).
			Post(fmt.Sprintf("%s/nodata/", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)

		var body map[string]interface{}
		json.Unmarshal([]byte(resp.String()), &body)

		if _, ok := body["id"]; ok {
			nid = int(body["id"].(float64))
		}
	})

	if nid > 0 {
		Convey("Delete nodata config", t, func() {
			resp, _ := rt.R().Delete(fmt.Sprintf("%s/nodata/%d", api_host, nid))
			So(resp.StatusCode(), ShouldEqual, 200)
		})
	}
}

func TestUser(t *testing.T) {

	Convey("Get User Login Failed", t, func() {
		rt := resty.New()
		resp, _ := rt.R().
			SetQueryParam("name", test_user_name).
			SetQueryParam("password", "willnotpass").
			Post(fmt.Sprintf("%s/user/login", api_host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Get User Login Success", t, func() {
		rt := resty.New()
		resp, _ := rt.R().SetQueryParam("name", test_user_name).
			SetQueryParam("password", test_user_password).
			Post(fmt.Sprintf("%s/user/login", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Session checking success", t, func() {
		apitoken, _ := get_session_token()

		rt := resty.New()
		rt.SetHeader("Apitoken", apitoken)
		resp, err := rt.R().Get(fmt.Sprintf("%s/user/auth_session", api_host))
		if err != nil {
			log.Error(err.Error())
		}
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Session checking failed", t, func() {
		invalid_apitoken := `{"name":"user-not-exists", "sig":"xxxx"}`
		rt := resty.New()
		rt.SetHeader("Apitoken", invalid_apitoken)
		resp, err := rt.R().Get(fmt.Sprintf("%s/user/auth_session", api_host))
		if err != nil {
			log.Error(err.Error())
		}
		So(resp.StatusCode(), ShouldEqual, 401)
	})

	Convey("Test Logout Session", t, func() {
		apitoken, _ := get_session_token()
		rt := resty.New()
		rt.SetHeader("Apitoken", apitoken)
		resp, err := rt.R().Get(fmt.Sprintf("%s/user/logout", api_host))
		if err != nil {
			log.Error(err.Error())
		}
		So(resp.StatusCode(), ShouldEqual, 200)
	})
}

func TestGraph(t *testing.T) {
	apitoken, _ := get_session_token()
	rt := resty.New()
	rt.SetHeader("Apitoken", apitoken)

	Convey("Get Endpoint Failed", t, func() {
		resp, _ := rt.R().Get(fmt.Sprintf("%s/graph/endpoint", api_host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Get Endpoint without login session", t, func() {
		resp, _ := resty.R().Get(fmt.Sprintf("%s/graph/endpoint", api_host))
		So(resp.StatusCode(), ShouldEqual, 401)
	})

	Convey("Get Endpoint List", t, func() {
		resp, _ := rt.R().SetQueryParam("q", "a.+").Get(fmt.Sprintf("%s/graph/endpoint", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Get Counter Failed", t, func() {
		resp, _ := rt.R().Get(fmt.Sprintf("%s/graph/endpoint_counter", api_host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Get Counter List", t, func() {
		resp, _ := rt.R().SetQueryParam("eid", "6,7").SetQueryParam("metricQuery", "disk.+").Get(fmt.Sprintf("%s/graph/endpoint_counter", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Delete counter", t, func() {
		resp, _ := rt.R().SetHeader("Content-Type", "application/json").
			SetBody(`{"endpoints":["laiwei-aggregator-1"], "counters":["agent.alive.percent/name=xx"]}`).
			Delete(fmt.Sprintf("%s/graph/counter", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Delete endpoint", t, func() {
		resp, _ := rt.R().SetHeader("Content-Type", "application/json").
			SetBody(`["0.0.0.0"]`).
			Delete(fmt.Sprintf("%s/graph/endpoint", api_host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

}

func TestTeam(t *testing.T) {
	apitoken, _ := get_session_token()
	test_team_name := "api-test-team1"
	test_team_id := 0

	Convey("Create Team Failed", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"api.test-resume": "i'm descript"}`).
			Post(fmt.Sprintf("%s/team", api_host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Create Team Scuessed", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", apitoken)
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"team_name": "%s","resume": "i'm descript", "users": [1]}`, test_team_name)).
			Post(fmt.Sprintf("%s/team", api_host))
		log.Debug(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	Convey("Get A Team by Name", t, func() {
		rt := resty.New()
		rt.SetHeader("Apitoken", apitoken)
		resp, _ := rt.R().
			Get(fmt.Sprintf("%s/team/name/%s", api_host, test_team_name))
		log.Debugf("reponsed: %v, team_id: %v", resp.String(), test_team_name)
		So(resp.StatusCode(), ShouldEqual, 200)

		var j map[string]interface{}
		json.Unmarshal([]byte(resp.String()), &j)

		if _, ok := j["id"]; ok {
			test_team_id = int(j["id"].(float64))
		}
	})

	if test_team_id > 0 {
		Convey("Delete A Team", t, func() {
			rt := resty.New()
			rt.SetHeader("Apitoken", apitoken)
			resp, _ := rt.R().
				Delete(fmt.Sprintf("%s/team/%d", api_host, test_team_id))
			log.Debugf("reponsed: %v, team_id: %v", resp.String(), test_team_id)
			So(resp.StatusCode(), ShouldEqual, 200)
		})
	}
}
