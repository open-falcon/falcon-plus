package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tidwall/gjson"
)

/*  convered routes test
 *	u.GET("/auth_session", AuthSession)
 *	authapi.GET("/current", UserInfo)
 *  authapi.GET("/u/:uid", GetUser)
 *  authapi.GET("/users", UserList)
 */
func TestUserSession(t *testing.T) {
	routes := SetUpGin()
	Convey("user session", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("session auth", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/user/auth_session", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "vaild")
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
		Convey("session auth faild", func() {
			w, r = NewTestContext("GET", "/api/v1/user/auth_session", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldNotContainSubstring, "vaild")
			So(w.Code, ShouldEqual, 401)
			CleanSession(r)
		})
		Convey("get current user info", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/user/current", nil)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "name\":\"testuser92")
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
	})
}

func TestUserGetInfo(t *testing.T) {
	routes := SetUpGin()
	Convey("get user info", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("get user info by id", func() {
			w, r = NewTestContext("GET", "/api/v1/user/u/3", nil)
			SetDefaultAdminSession(r)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "id\":3")
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
	})
	Convey("get user list", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("user list", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/user/users", nil)
			routes.ServeHTTP(w, r)
			checkR := gjson.Get(w.Body.String(), "#.name")
			So(len(checkR.Array()), ShouldBeGreaterThan, 2)
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
		Convey("user list with page & limit", func() {
			w, r = NewTestContextWithDefaultSession("GET", "/api/v1/user/users?page=1&limit=2", nil)
			routes.ServeHTTP(w, r)
			checkR := gjson.Get(w.Body.String(), "#.name")
			So(len(checkR.Array()), ShouldEqual, 2)
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
	})
}
