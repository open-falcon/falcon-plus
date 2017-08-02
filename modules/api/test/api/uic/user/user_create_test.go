package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

/*  convered routes test
 *	u.POST("/create", CreateUser)
 */
// only can run once, after this should need reconver db
func TestUserCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("test create user", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create a normal user", func() {
			postb := map[string]interface{}{
				"name":     "testuser93",
				"cnname":   "testuser93",
				"password": "testuser93",
				"email":    "testuser93@open.com",
				"phone":    "000-000-000-0000",
				"im":       "000000000000",
				"qq":       "000000000000",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("POST", "/api/v1/user/create", &b)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldContainSubstring, "{\"name\":")
			So(w.Code, ShouldEqual, 200)
		})
		// Convey("create a root user", func() {
		// 	rootTokens := ""
		// 	postb := map[string]interface{}{
		// 		"name":     "root",
		// 		"cnname":   "管理員",
		// 		"password": "root1234",
		// 		"email":    "root@open.com",
		// 		"phone":    "000-000-000-0000",
		// 		"im":       "000000000000",
		// 		"qq":       "000000000000",
		// 	}
		// 	b, _ := json.Marshal(postb)
		// 	w, r = NewTestContext("POST", "/api/v1/user/create", &b)
		// 	routes.ServeHTTP(w, r)
		// 	rootTokens = w.Body.String()
		// 	So(rootTokens, ShouldContainSubstring, "{\"name\":")
		// 	So(w.Code, ShouldEqual, 200)
		// 	Convey("create a normal user with admin privileges", func() {
		// 		postb := map[string]interface{}{
		// 			"name":     "testuser90",
		// 			"cnname":   "testuser90",
		// 			"password": "testuser90",
		// 			"email":    "testuser90@open.com",
		// 			"phone":    "000-000-000-0000",
		// 			"im":       "000000000000",
		// 			"qq":       "000000000000",
		// 		}
		// 		b, _ := json.Marshal(postb)
		// 		w, r = NewTestContext("POST", "/api/v1/user/create", &b)
		// 		sname, ssig := hutil.ParseCookieFromResp(rootTokens)
		// 		r = SetSessionWith(r, sname, ssig)
		// 		viper.Set("signup_disable", true)
		// 		routes.ServeHTTP(w, r)
		// 		rootTokens = w.Body.String()
		// 		So(rootTokens, ShouldContainSubstring, "{\"name\":")
		// 		So(w.Code, ShouldEqual, 200)
		// 		CleanSession(r)
		// 	})
		// })
	})
}

func TestSingUpDisableCreateUser(t *testing.T) {
	routes := SetUpGin()
	Convey("test create user with sign up disabled", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create a user91 user", func() {
			viper.Set("signup_disable", true)
			postb := map[string]interface{}{
				"name":     "testuser91",
				"cnname":   "testuser91",
				"password": "testuser91",
				"email":    "testuser91@open.com",
				"phone":    "000-000-000-0000",
				"im":       "000000000000",
				"qq":       "000000000000",
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContext("POST", "/api/v1/user/create", &b)
			CleanSession(r)
			routes.ServeHTTP(w, r)
			So(w.Body.String(), ShouldContainSubstring, "sign up is not enabled")
			So(w.Code, ShouldEqual, 400)
			log.Info(w.Body.String())
		})
	})
}

// func TestOneUser(t *testing.T) {
// 	routes := SetUpGin()
// 	Convey("test create user with sign up disabled2", t, func() {
// 		var (
// 			w *httptest.ResponseRecorder
// 			r *http.Request
// 		)
// 		Convey("create a test03 user", func() {
// 			postb := map[string]interface{}{
// 				"name":     "test03",
// 				"cnname":   "test03",
// 				"password": "test03",
// 				"email":    "test03@open.com",
// 				"phone":    "000-000-000-0000",
// 				"im":       "000000000000",
// 				"qq":       "000000000000",
// 			}
// 			b, _ := json.Marshal(postb)
// 			w, r = NewTestContext("POST", "/api/v1/user/create", &b)
// 			CleanSession(r)
// 			routes.ServeHTTP(w, r)
// 			So(w.Body.String(), ShouldNotContainSubstring, "sign up is not enabled")
// 			So(w.Code, ShouldEqual, 200)
// 			log.Info(w.Body.String())
// 		})
// 	})
// }
