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
 *	authapi_team.POST("/team", CreateTeam)
 */

func TestTeamCreate(t *testing.T) {
	routes := SetUpGin()
	Convey("create a new team", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create a new team ok", func() {
			postb := map[string]interface{}{
				"team_name": "team_X",
				"resume":    "this is resumeA",
				"users":     []int{1, 2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "team created")
			So(w.Code, ShouldEqual, 200)
		})
		Convey("create a new team faild", func() {
			postb := map[string]interface{}{
				"resume":  "this is resume3",
				"userIDs": []int{2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 400)
		})
	})
}
