package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tidwall/gjson"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
)

/*  convered routes test
*	authapi_team.PUT("/team", UpdateTeam)
 */

func TestTeamUpdate(t *testing.T) {
	routes := SetUpGin()
	Convey("create a new team", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("update a new team ok", func() {
			postb := map[string]interface{}{
				"team_id":   4,
				"team_name": "team_D",
				"resume":    "this is resumeD",
				"users":     []int{1},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			checkR := gjson.Get(respBody, "users.#.id")
			So(respBody, ShouldContainSubstring, "team_D")
			So(len(checkR.Array()), ShouldEqual, 1)
			So(w.Code, ShouldEqual, 200)
		})
		Convey("update a new team faild (with missing id)", func() {
			postb := map[string]interface{}{
				"team_name": "team_E",
				"resume":    "this is resume3",
				"userIDs":   []int{2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("PUT", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, 400)
		})
	})
}
