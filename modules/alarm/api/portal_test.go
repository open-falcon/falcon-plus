package api

import (
	"testing"

	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	g.ParseConfig("../cfg.example.json")
}

func TestPortalAPI(t *testing.T) {
	Convey("Get action from api failed", t, func() {
		r := CurlAction(1)
		So(r.Id, ShouldEqual, 1)
	})

}
