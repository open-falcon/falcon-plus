package thelp

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/controller"
	"github.com/open-falcon/falcon-plus/modules/api/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var routes *gin.Engine

func SetUpGin() *gin.Engine {
	if routes != nil {
		return routes
	} else {
		confPath := flag.String("conf", "test_cfg", "set test configure file's name")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/")
		viper.AddConfigPath("../../../")
		viper.AddConfigPath("../../../../")
		viper.SetConfigName(*confPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Error(err.Error())
		}
		gin.SetMode(gin.TestMode)
		log.SetLevel(log.DebugLevel)
		config.InitDB(viper.GetBool("db.db_debug"))
		//test with default set of db
		routes := gin.Default()
		routes = controller.StartGin(":9898", routes, true)
		return routes
	}
}
