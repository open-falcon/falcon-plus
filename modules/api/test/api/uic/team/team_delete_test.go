package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/open-falcon/falcon-plus/modules/api/test_utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tidwall/gjson"
)

/*  convered routes test
*	authapi_team.DELETE("/team", DeleteTeam)
 */

func TestDeleteTeam(t *testing.T) {
	routes := SetUpGin()
	Convey("delete a team", t, func() {
		var (
			w *httptest.ResponseRecorder
			r *http.Request
		)
		Convey("create a team", func() {
			// create a new team for delete
			postb := map[string]interface{}{
				"team_name": "team_test3",
				"resume":    "this is resumeA",
				"users":     []int{1, 2, 3},
			}
			b, _ := json.Marshal(postb)
			w, r = NewTestContextWithDefaultSession("POST", "/api/v1/team", &b)
			routes.ServeHTTP(w, r)
			respBody := w.Body.String()
			So(respBody, ShouldContainSubstring, "team created")
			So(w.Code, ShouldEqual, 200)
			Tid := gjson.Get(w.Body.String(), "team.id")
			Convey("delete a team ok", func() {
				team_id := Tid.String()
				w, r = NewTestContextWithDefaultSession("DELETE", "/api/v1/team/"+team_id, nil)
				routes.ServeHTTP(w, r)
				respBody := w.Body.String()
				So(respBody, ShouldContainSubstring, "deleted. Affect row")
				So(w.Code, ShouldEqual, 200)
			})
		})
	})
}
