package dashboard_screen

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	authapi := r.Group("/api/v1/dashboard")
	authapi.Use(utils.AuthSessionMidd)
	authapi.POST("/screen", ScreenCreate)
	authapi.GET("/screen/:screen_id", ScreenGet)
	authapi.GET("/screens/pid/:pid", ScreenGetsByPid)
	authapi.GET("/screens", ScreenGetsAll)
	authapi.DELETE("/screen/:screen_id", ScreenDelete)
	authapi.PUT("/screen/:screen_id", ScreenUpdate)
}
