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
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/masato25/resty"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	cutils "github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	cfg "github.com/open-falcon/falcon-plus/modules/api/config"
)

var (
	api_v1             = ""
	test_user_name     = "apitest-user1"
	test_user_password = "password"
	test_team_name     = "apitest-team1"
	root_user_name     = "root"
	root_user_password = "rootpass"
)

func init() {
	cfg_file := os.Getenv("API_TEST_CFG")
	if cfg_file == "" {
		cfg_file = "./cfg.example"
	}
	viper.SetConfigName(cfg_file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.InitLog(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}

	db_user := os.Getenv("DB_USER")
	if db_user == "" {
		db_user = "root"
	}

	db_passwd := os.Getenv("DB_PASSWORD")

	db_host := os.Getenv("DB_HOST")
	if db_host == "" {
		db_host = "127.0.0.1"
	}

	db_port := os.Getenv("DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}

	db_names := []string{"falcon_portal", "graph", "uic", "dashboard", "alarms"}
	for _, dbn := range db_names {
		viper.Set(fmt.Sprintf("db.%s", dbn), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			db_user, db_passwd, db_host, db_port, dbn))
	}

	err = cfg.InitDB(viper.GetBool("db.db_bug"), viper.GetViper())
	if err != nil {
		log.Fatalf("db conn failed with error %s", err.Error())
	}

	api_port := os.Getenv("API_PORT")
	if api_port == "" {
		api_port = strings.TrimLeft(viper.GetString("web_port"), ":")
	}
	api_host := os.Getenv("API_HOST")
	if api_host == "" {
		api_host = "127.0.0.1"
	}
	api_v1 = fmt.Sprintf("http://%s:%s/api/v1", api_host, api_port)

	init_testing_data()
}

func init_testing_data() {
	password := utils.HashIt(test_user_password)
	user := uic.User{
		Name:   test_user_name,
		Passwd: password,
		Cnname: test_user_name,
		Email:  test_user_name + "@test.com",
		Phone:  "1234567890",
		IM:     "hellotest",
		QQ:     "3800000",
	}

	db := cfg.Con()
	if db.Uic.Table("user").Where("name = ?", test_user_name).First(&uic.User{}).RecordNotFound() {
		if err := db.Uic.Table("user").Create(&user).Error; err != nil {
			log.Fatal(err)
		}
		log.Info("create_user:", test_user_name)
	}

	db.Uic.Table("user").Where("name = ?", "root").Delete(&uic.User{})
	db.Uic.Table("team").Where("name = ?", test_team_name).Delete(&uic.Team{})
}

func get_session_token() (string, error) {
	rr := map[string]interface{}{}
	resp, _ := resty.R().
		SetQueryParam("name", root_user_name).
		SetQueryParam("password", root_user_password).
		SetResult(&rr).
		Post(fmt.Sprintf("%s/user/login", api_v1))

	if resp.StatusCode() != 200 {
		return "", errors.New(resp.String())
	}

	api_token := fmt.Sprintf(`{"name": "%v", "sig": "%v"}`, rr["name"], rr["sig"])
	return api_token, nil
}

func TestUser(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}
	var api_token string = ""

	Convey("Create root user: POST /user/create", t, func() {
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]string{
				"name":     root_user_name,
				"password": root_user_password,
				"email":    "root@test.com",
				"cnname":   "cnroot",
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/create", api_v1))

		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["sig"], ShouldNotBeBlank)
		api_token = resp.String()
	})

	Convey("Get user info by name: GET /user/name/:user", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["role"], ShouldEqual, 2)
		So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
	})
	root_user_id := (*rr)["id"]

	Convey("Get user info by id: GET /user/u/:uid", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/u/%v", api_v1, root_user_id))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, root_user_name)
	})

	Convey("Update current user: PUT /user/update", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(map[string]string{
				"cnname": "cnroot2",
				"email":  "root2@test.com",
				"phone":  "18000000000",
			}).
			SetResult(rr).
			Put(fmt.Sprintf("%s/user/update", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "updated")

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			resp, _ := resty.R().
				SetHeader("Apitoken", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["cnname"], ShouldEqual, "cnroot2")
		})
	})

	Convey("Change password: PUT /user/cgpasswd", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(map[string]string{
				"old_password": root_user_password,
				"new_password": root_user_password,
			}).
			SetResult(rr).
			Put(fmt.Sprintf("%s/user/cgpasswd", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "updated")
	})

	Convey("Get user list: GET /user/users", t, func() {
		r := []map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(&r).
			Get(fmt.Sprintf("%s/user/users", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(r, ShouldNotBeEmpty)
		So(r[0]["name"], ShouldNotBeBlank)
	})

	Convey("Get current user: POST /user/current", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/current", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, root_user_name)
	})

	Convey("Login user: POST /user/login", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetQueryParam("name", root_user_name).
			SetQueryParam("password", root_user_password).
			SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, root_user_name)
		So((*rr)["sig"], ShouldNotBeBlank)
		So((*rr)["admin"], ShouldBeTrue)
		api_token = fmt.Sprintf(`{"name": "%v", "sig": "%v"}`, (*rr)["name"], (*rr)["sig"])
	})

	Convey("Auth user by session: GET /user/auth_session", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/auth_session", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "valid")
	})

	Convey("Logout user: GET /user/logout", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/logout", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "successful")
	})
}

func TestAdmin(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}
	var api_token string = ""

	Convey("Login as root", t, func() {
		resp, _ := resty.R().
			SetQueryParam("name", root_user_name).SetQueryParam("password", root_user_password).SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, root_user_name)
		So((*rr)["sig"], ShouldNotBeBlank)
		So((*rr)["admin"], ShouldBeTrue)
	})
	api_token = fmt.Sprintf(`{"name": "%v", "sig": "%v"}`, (*rr)["name"], (*rr)["sig"])

	Convey("Get user info by name: GET /user/name/:user", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
	})
	test_user_id := (*rr)["id"]

	Convey("Change user role: PUT /admin/change_user_role", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"user_id": %v,"admin": "yes"}`, test_user_id)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_role", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "sccuessful")

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			resp, _ := resty.R().
				SetHeader("Apitoken", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["role"], ShouldEqual, 1)
		})
	})

	Convey("Change user passwd: PUT /admin/change_user_passwd", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"user_id": %v,"password": "%s"}`, test_user_id, test_user_password)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_passwd", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "updated")
	})

	Convey("Change user profile: PUT /admin/change_user_profile", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"user_id": %v,"cnname": "%s", "email": "%s"}`,
				test_user_id, test_user_name, "test_user1@test.com")).
			SetResult(rr).
			Put(fmt.Sprintf("%s/admin/change_user_profile", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "updated")

		Convey("Get user info by name: GET /user/name/:user", func() {
			*rr = map[string]interface{}{}
			resp, _ := resty.R().
				SetHeader("Apitoken", api_token).
				SetResult(rr).
				Get(fmt.Sprintf("%s/user/name/%s", api_v1, test_user_name))
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["email"], ShouldEqual, "test_user1@test.com")
		})
	})

	Convey("Admin login user: POST /admin/login", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(map[string]string{
				"name": test_user_name,
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/admin/login", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, test_user_name)
	})

	Convey("Delete user: DELETE /admin/delete_user", t, func() {
	})
}

func TestTeam(t *testing.T) {
	var rr *map[string]interface{} = &map[string]interface{}{}

	Convey("Login as root", t, func() {
		resp, _ := resty.R().
			SetQueryParam("name", root_user_name).SetQueryParam("password", root_user_password).SetResult(rr).
			Post(fmt.Sprintf("%s/user/login", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, root_user_name)
		So((*rr)["sig"], ShouldNotBeBlank)
		So((*rr)["admin"], ShouldBeTrue)
	})
	api_token := fmt.Sprintf(`{"name": "%v", "sig": "%v"}`, (*rr)["name"], (*rr)["sig"])

	Convey("Get user info by name: GET /user/name/:user", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/name/%s", api_v1, root_user_name))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["role"], ShouldEqual, 2)
		So((*rr)["id"], ShouldBeGreaterThanOrEqualTo, 0)
	})
	root_user_id := (*rr)["id"]

	Convey("Create team: POST /team", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(fmt.Sprintf(`{"team_name": "%s","resume": "i'm descript", "users": [1]}`, test_team_name)).
			SetResult(rr).
			Post(fmt.Sprintf("%s/team", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "created")
	})

	Convey("Get team by name: GET /team/name/:name", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().SetHeader("Apitoken", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/team/name/%s", api_v1, test_team_name))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, test_team_name)
		So((*rr)["users"], ShouldNotBeEmpty)
		So((*rr)["id"], ShouldBeGreaterThan, 0)
	})
	test_team_id := (*rr)["id"]

	Convey("Get team by id: GET /team/t/:tid", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().SetHeader("Apitoken", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/team/t/%v", api_v1, test_team_id))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["name"], ShouldEqual, test_team_name)
		So((*rr)["users"], ShouldNotBeEmpty)
		So((*rr)["id"], ShouldEqual, test_team_id)
	})

	Convey("Update team by id: PUT /team", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(fmt.Sprintf(`{"team_id": %v,"resume": "descript2", "name":"%v", "users": [1]}`,
				test_team_id, test_team_name)).
			SetResult(rr).
			Put(fmt.Sprintf("%s/team", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "updated")

		Convey("Get team by name: GET /team/name/:name", func() {
			*rr = map[string]interface{}{}
			resp, _ := resty.R().SetHeader("Apitoken", api_token).SetResult(rr).
				Get(fmt.Sprintf("%s/team/name/%s", api_v1, test_team_name))
			So(resp.StatusCode(), ShouldEqual, 200)
			So(*rr, ShouldNotBeEmpty)
			So((*rr)["resume"], ShouldEqual, "descript2")
		})
	})

	Convey("Add users to team: POST /team/user", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Apitoken", api_token).
			SetBody(map[string]interface{}{
				"team_id": test_team_id,
				"users":   []string{root_user_name},
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/team/user", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "successful")
	})

	Convey("Get teams which user belong to: GET /user/u/:uid/teams", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().SetHeader("Apitoken", api_token).SetResult(rr).
			Get(fmt.Sprintf("%s/user/u/%v/teams", api_v1, root_user_id))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["teams"], ShouldNotBeEmpty)
	})

	Convey("Check user in teams or not: GET /user/u/:uid/in_teams", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetQueryParam("team_names", test_team_name).
			SetResult(rr).
			Get(fmt.Sprintf("%s/user/u/%v/in_teams", api_v1, root_user_id))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldEqual, "true")
	})

	Convey("Get team list: GET /team", t, func() {
		var r []map[string]interface{}
		resp, _ := resty.R().SetHeader("Apitoken", api_token).SetResult(&r).
			Get(fmt.Sprintf("%s/team", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(r, ShouldNotBeEmpty)
		So(r[0]["team"], ShouldNotBeEmpty)
		So(r[0]["users"], ShouldNotBeEmpty)
		So(r[0]["creator_name"], ShouldNotBeBlank)
	})

	Convey("Delete team by id: DELETE /team/:tid", t, func() {
		*rr = map[string]interface{}{}
		resp, _ := resty.R().
			SetHeader("Apitoken", api_token).
			SetResult(rr).
			Delete(fmt.Sprintf("%s/team/%v", api_v1, test_team_id))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
		So((*rr)["message"], ShouldContainSubstring, "deleted")
	})
}

func TestGraph(t *testing.T) {
	api_token, err := get_session_token()
	if err != nil {
		log.Fatal(err)
	}

	rc := resty.New()
	rc.SetHeader("Apitoken", api_token)
	var rr *[]map[string]interface{} = &[]map[string]interface{}{}

	Convey("Get endpoint list: GET /graph/endpoint", t, func() {
		resp, _ := rc.R().SetQueryParam("q", ".+").
			SetResult(rr).
			Get(fmt.Sprintf("%s/graph/endpoint", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(*rr), ShouldBeGreaterThanOrEqualTo, 0)

	})

	if len(*rr) == 0 {
		return
	}

	eid := (*rr)[0]["id"]
	endpoint := (*rr)[0]["endpoint"]
	Convey("Get counter list: GET /graph/endpoint_counter", t, func() {
		resp, _ := rc.R().
			SetQueryParam("eid", fmt.Sprintf("%v", eid)).
			SetQueryParam("metricQuery", ".+").
			SetQueryParam("limit", "1").
			SetResult(rr).
			Get(fmt.Sprintf("%s/graph/endpoint_counter", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
	})

	if len(*rr) == 0 {
		return
	}

	counter := (*rr)[0]["counter"]
	step := (*rr)[0]["step"]

	now := time.Now()
	start_ts := now.Add(time.Duration(-1) * time.Hour).Unix()
	end_ts := now.Unix()

	Convey("Query counter history: POST /graph/history", t, func() {
		resp, _ := rc.R().
			SetBody(map[string]interface{}{
				"step":       step,
				"consol_fun": "AVERAGE",
				"start_time": start_ts,
				"end_time":   end_ts,
				"hostnames":  []string{endpoint.(string)},
				"counters":   []string{counter.(string)},
			}).
			SetResult(rr).
			Post(fmt.Sprintf("%s/graph/history", api_v1))
		log.Info(resp)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(*rr, ShouldNotBeEmpty)
	})
}

func TestNodata(t *testing.T) {
	api_token, err := get_session_token()
	if err != nil {
		log.Fatal(err)
	}

	var rr *map[string]interface{} = &map[string]interface{}{}
	rc := resty.New()
	rc.SetHeader("Apitoken", api_token)

	var nid int = 0

	Convey("Create nodata config: POST /nodata", t, func() {
		nodata_name := fmt.Sprintf("api.testnodata-%s", cutils.RandString(8))
		resp, _ := rc.R().
			SetHeader("Content-Type", "application/json").
			SetBody(fmt.Sprintf(`{"tags": "", "step": 60, "obj_type": "host", "obj": "docker-agent",
				"name": "%s", "mock": -1, "metric": "api.test.metric", "dstype": "GAUGE"}`, nodata_name)).
			SetResult(rr).
			Post(fmt.Sprintf("%s/nodata/", api_v1))
		So(resp.StatusCode(), ShouldEqual, 200)

		if v, ok := (*rr)["id"]; ok {
			nid = int(v.(float64))
			Convey("Delete nodata config", func() {
				resp, _ := rc.R().Delete(fmt.Sprintf("%s/nodata/%d", api_v1, nid))
				So(resp.StatusCode(), ShouldEqual, 200)
			})
		}
	})
}
