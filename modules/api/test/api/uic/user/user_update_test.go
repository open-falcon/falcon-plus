package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
 *	authapi.PUT("/cgpasswd", ChangePassword)
 *	authapi.PUT("/update", UpdateUser)
 */
func TestUsetChngePWD(t *testing.T) {
	routes := SetUpGin()
	Convey("change user password", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("update user password ok", func() {
			postb := map[string]interface{}{
				"old_password": "test02",
				"new_password": "test02",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/user/cgpasswd", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "password updated")
			So(w.Code, ShouldEqual, 200)
		})
		Convey("update user password error", func() {
			postb := map[string]interface{}{
				"old_password": "testerror",
				"new_password": "test02",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/user/cgpasswd", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "\"error\"")
			So(w.Code, ShouldEqual, 400)
		})
	})
}

func TestUserUpdate(t *testing.T) {
	routes := SetUpGin()
	Convey("user update", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("update user info", func() {
			postb := map[string]interface{}{
				"cnname": "newcnname",
				"email":  "newemail@gmail.com",
				"phone":  "2222-2222-222",
				"im":     "2222222222222",
				"qq":     "2222222222222",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/user/update", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "user info updated")
			So(w.Code, ShouldEqual, 200)
		})
	})
}
