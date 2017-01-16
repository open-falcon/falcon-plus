package test

import (
	"fmt"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
	"github.com/masato25/resty"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCUser(t *testing.T) {
	viper.AddConfigPath("../../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	host := "http://localhost:3000/api/v1"
	Convey("Register User Failed", t, func() {
		rt := resty.New()
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"name": "test","password": "test"}`).
			Post(fmt.Sprintf("%s/user/create", host))
		So(resp.StatusCode(), ShouldEqual, 400)
	})

	Convey("Register User Scuessed", t, func() {
		rt := resty.New()
		resp, _ := rt.R().
			SetHeader("Content-Type", "application/json").
			SetBody(`{"name": "owltester","password": "mypassword", "cnname": "翱鶚", "email": "root123@cepave.com", "im": "44955834958", "phone": "99999999999", "qq": "904394234239"}`).
			Post(fmt.Sprintf("%s/user/create", host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})
	Convey("Update User", t, func() {
		Convey("Update User Scuessed", func() {
			cname := "test1"
			csig := "d4f71cba377911e699d60242ac110010"
			Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
			rt := resty.New()
			rt.SetHeader("Apitoken", Apitoken)
			resp, _ := rt.R().
				SetHeader("Content-Type", "application/json").
				SetBody(`{"name": "test1","password": "test1", "cnname": "翱鶚Test", "email": "root123@cepave.com", "im": "44955834958", "phone": "99999999999", "qq": "904394234239"}`).
				Put(fmt.Sprintf("%s/user/update", host))
			So(resp.StatusCode(), ShouldEqual, 200)

			Convey("Update User Password", func() {
				cname := "test1"
				csig := "d4f71cba377911e699d60242ac110010"
				Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
				rt := resty.New()
				rt.SetHeader("Apitoken", Apitoken)
				resp, _ := rt.R().
					SetHeader("Content-Type", "application/json").
					SetBody(`{"new_password": "test1", "old_password": "test1"}`).
					Put(fmt.Sprintf("%s/user/cgpasswd", host))
				So(resp.StatusCode(), ShouldEqual, 200)
			})
		})
	})

	Convey("Get User List", t, func() {
		cname := "test1"
		csig := "d4f71cba377911e699d60242ac110010"
		Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
		rt := resty.New()
		rt.SetHeader("Apitoken", Apitoken)
		resp, _ := rt.R().Get(fmt.Sprintf("%s/user/users", host))
		So(resp.StatusCode(), ShouldEqual, 200)
	})

	// Convey("Change User Role as Admin", t, func() {
	// 	cname := "root"
	// 	csig := "dd81ea033c2d11e6a95d0242ac11000c"
	// 	Apitoken := fmt.Sprintf(`{"name": "%s", "sig": "%s"}`, cname, csig)
	// 	rt := resty.New()
	// 	rt.SetHeader("Apitoken", Apitoken)
	// 	resp, _ := rt.R().
	// 		SetHeader("Content-Type", "application/json").
	// 		SetBody(`{"user_id": 4, "admin": "yes"}`).
	// 		Put(fmt.Sprintf("%s/user/cgrole", host))
	// 	So(resp.StatusCode(), ShouldEqual, 200)
	// })
	//
	// Convey("Change User Role as normal user", t, func() {
	// 	cname := "root"
	// 	csig := "dd81ea033c2d11e6a95d0242ac11000c"
	// 	rt := resty.New()
	// 	rt.SetHeader("Apitoken", Apitoken)
	// 	resp, _ := rt.R().
	// 		SetHeader("Content-Type", "application/json").
	// 		SetBody(`{"user_id": 4, "admin": "no"}`).
	// 		Put(fmt.Sprintf("%s/user/cgrole", host))
	// 	So(resp.StatusCode(), ShouldEqual, 200)
	// })
}
