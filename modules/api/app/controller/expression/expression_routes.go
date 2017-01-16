package expression

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/masato25/owl_backend/app/utils"
	"github.com/masato25/owl_backend/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest
const expecstatus = http.StatusExpectationFailed

func Routes(r *gin.Engine) {
	db = config.Con()
	expr := r.Group("/api/v1/expression")
	expr.Use(utils.AuthSessionMidd)
	expr.GET("", GetExpressionList)
	expr.GET("/:eid", GetExpression)
	expr.POST("", CreateExrpession)
	expr.PUT("", UpdateExrpession)
	expr.DELETE("/:eid", DeleteExpression)
}
