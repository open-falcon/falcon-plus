package sdk

import (
	"testing"

	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	g.ParseConfig("../cfg.example.json")
}

func TestSDK(t *testing.T) {
	Convey("get hostnames by id", t, func() {
		r, err := HostnamesByID(1)
		t.Log(r, err)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeGreaterThanOrEqualTo, 0)
	})

	Convey("query last points", t, func() {
		r, err := QueryLastPoints([]string{"laiweiofficemac"}, []string{"agent.alive"})
		t.Log(r, err)
		So(err, ShouldBeNil)
		So(len(r), ShouldBeGreaterThanOrEqualTo, 0)
		for _, x := range r {
			t.Log(x)
		}
	})
}
