package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	hutil "github.com/open-falcon/falcon-plus/modules/api/test_utils/util"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
 *	u.POST("/login", Login)
 *	u.GET("/logout", Logout)
 *  u.GET("/auth_session", AuthSession)
 */
func TestUserLogin(t *testing.T) {
	routes := SetUpGin()
	Convey("login user", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("test successful - login user & session auth & auth session &logout user", func() {
			postb := map[string]interface{}{
				"name":     "testuser99",
				"password": "testuser99",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("POST", "/api/v1/user/login", &b)
			routes.ServeHTTP(w, r)
			loginTokens := w.Body.String()
			So(loginTokens, ShouldContainSubstring, "\"sig\":")
			So(w.Code, ShouldEqual, 200)
			sname, ssig := hutil.ParseCookieFromResp(loginTokens)
			Convey("test auth session", func() {
				w, r := NewTestContext("GET", "/api/v1/user/auth_session", nil)
				SetSessionWith(r, sname, ssig)
				routes.ServeHTTP(w, r)
				So(w.Body.String(), ShouldContainSubstring, "vaild")
				So(w.Code, ShouldEqual, 200)
				CleanSession(r)
				Convey("test logout user", func() {
					w, r := NewTestContext("GET", "/api/v1/user/logout", nil)
					SetSessionWith(r, sname, ssig)
					routes.ServeHTTP(w, r)
					So(w.Body.String(), ShouldContainSubstring, "logout successful")
					So(w.Code, ShouldEqual, 200)
					CleanSession(r)
				})
			})
		})
		Convey("test login user faild", func() {
			postb := map[string]interface{}{
				"name":     "testuser99",
				"password": "0000",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("POST", "/api/v1/user/login", &b)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldNotContainSubstring, "\"sig\":")
			So(w.Code, ShouldEqual, 400)
		})
		Convey("test auth session faild", func() {
			w, r := NewTestContext("GET", "/api/v1/user/auth_session", nil)
			SetSessionWith(r, "testuser99", "0000")
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldContainSubstring, "not found")
			So(w.Code, ShouldEqual, 401)
			CleanSession(r)
		})
	})
}
