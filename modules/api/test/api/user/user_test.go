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

func TestUser(t *testing.T) {
	viper.AddConfigPath("../../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:3000/api/v1"
	Convey("Get User Login Failed", t, func() {
		rt := resty.New()
		resp, _ := rt.R().SetQueryParam("name", "gg123").
			SetQueryParam("name", "root").
			SetQueryParam("password", "willnotpass").
			Post(fmt.Sprintf("%s/user/login", host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("Get User Login Success", t, func() {
		rt := resty.New()
		resp, _ := rt.R().SetQueryParam("name", "test2").
			SetQueryParam("password", "test2").
			Post(fmt.Sprintf("%s/user/login", host))
		log.Info("result: ", resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})
	Convey("Test Logout Session", t, func() {
		rt := resty.New()
		resp, _ := rt.R().SetQueryParam("name", "test2").
			SetQueryParam("password", "test2").
			Post(fmt.Sprintf("%s/user/login", host))
		jss, err := gojq.NewStringQuery(resp.String())
		sig, err := jss.Query("sig")
		Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, "test2", sig.(string))
		rt = resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, err = rt.R().Get(fmt.Sprintf("%s/user/logout", host))
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})
	Convey("Session checking success", t, func() {
		cname := "test1"
		csig := "d4f71cba377911e699d60242ac110010"
		Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, err := rt.R().Get(fmt.Sprintf("%s/user/auth_session", host))
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(resp.String())
		So(resp.StatusCode(), ShouldEqual, 200)
	})
	Convey("Session checking failed", t, func() {
		cname := "testtest"
		csig := "9a84ae1d377911e699d60242ac110010"
		Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, err := rt.R().Get(fmt.Sprintf("%s/user/auth_session", host))
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(resp.String())
		// log.Info(resp.Body.ToString())
		So(resp.StatusCode(), ShouldEqual, 401)
	})
}
