package test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

func TestHash(t *testing.T) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("cfg_test")
	viper.ReadInConfig()
	log.SetLevel(log.DebugLevel)
	Convey("Test Hash method", t, func() {
		val := utils.HashIt("test2")
		So(val, ShouldEqual, "c0fc7c3e09f7efc71567b453ec5b9cd2")
	})
}
