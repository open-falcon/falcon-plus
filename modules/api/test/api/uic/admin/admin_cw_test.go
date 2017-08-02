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
*	adminapi.PUT("/change_user_role", ChangeRuleOfUser)
*	adminapi.PUT("/change_user_passwd", ChangeRuleOfUser)
 */

func TestAdminChange(t *testing.T) {
	routes := SetUpGin()
	Convey("admin action test", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("change user role ok", func() {
			postb := map[string]interface{}{
				"user_id": 3,
				"admin":   "yes",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("PUT", "/api/v1/admin/change_user_role", &b)
			r = SetDefaultAdminSession(r)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "update sccuessful")
			So(w.Code, ShouldEqual, 200)
			Convey("reset user to normal permission", func() {
				postb := map[string]interface{}{
					"user_id": 3,
					"admin":   "no",
				}
				b, _ := json.Marshal(postb)
				w, r = NewTestContext("PUT", "/api/v1/admin/change_user_role", &b)
				r = SetDefaultAdminSession(r)
				routes.ServeHTTP(w, r)
				respBody := w.Body.String()
				So(respBody, ShouldContainSubstring, "update sccuessful")
				So(w.Code, ShouldEqual, 200)
			})
			CleanSession(r)
		})
		Convey("change user role failed", func() {
			postb := map[string]interface{}{
				"user_id": 3,
				"admin":   "yes",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/admin/change_user_role", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "don't have permission")
			So(w.Code, ShouldEqual, 400)
			CleanSession(r)
		})
		Convey("change user password ok", func() {
			postb := map[string]interface{}{
				"user_id":  3,
				"password": "testuser92",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("PUT", "/api/v1/admin/change_user_passwd", &b)
			r = SetDefaultAdminSession(r)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "password updated!")
			So(w.Code, ShouldEqual, 200)
			CleanSession(r)
		})
		Convey("change user password failed", func() {
			postb := map[string]interface{}{
				"user_id":  3,
				"password": "testuser92",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/admin/change_user_passwd", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "don't have permission")
			So(w.Code, ShouldEqual, 400)
			CleanSession(r)
		})
	})
}

func TestAdminDeleteUser(t *testing.T) {
	routes := SetUpGin()
	var (
		w *httptest.ResponseRecorder
		r *http.Request
	)
	Convey("create a new user", t, func() {
		postb := map[string]interface{}{
			"name":     "testuserd1",
			"cnname":   "testuserd1",
			"password": "testuserd1",
			"email":    "testuserd1@open.com",
			"phone":    "000-000-000-0000",
			"im":       "000000000000",
			"qq":       "000000000000",
		}
		b, _ := json.Marshal(postb)
		w, r = NewTestContext("POST", "/api/v1/user/create", &b)
		routes.ServeHTTP(w, r)
		So(w.Body.String(), ShouldContainSubstring, "\"name\":")
		So(w.Code, ShouldEqual, 200)
		var bindJosn map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &bindJosn)
		Convey("delete user ok", func() {
			postb := map[string]interface{}{
				"user_id": bindJosn["id"].(float64),
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("DELETE", "/api/v1/admin/delete_user", &b)
			r = SetDefaultAdminSession(r)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldContainSubstring, "has been delete")
			So(w.Code, ShouldEqual, 200)
		})
		Convey("delete user no permission", func() {
			postb := map[string]interface{}{
				"user_id": bindJosn["id"].(float64),
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("DELETE", "/api/v1/admin/delete_user", &b)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldContainSubstring, "don't have permission!")
			So(w.Code, ShouldEqual, 400)
		})
	})
}
