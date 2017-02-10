package dashboard_screen

import (
	"net/http"

	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
	"gopkg.in/gin-gonic/gin.v1"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/dashboard/screen")
	authapi.Use(utils.AuthSessionMidd)
	authapi.POST("", ScreenCreate)
	authapi.GET("/:screen_id", ScreenGet)
	authapi.GET("/pid/:pid", ScreenGetsByPid)
}
