package http

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/falcon2p8s/g"
)

func configCommonRoutes() {
	router.GET("/version", func(c *gin.Context) {
		JSONR(c, 200, gin.H{"value": g.VersionMsg()})
	})

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "ok")
	})

	router.POST("/config/reload", func(c *gin.Context) {
		g.ParseConfig(g.ConfigFile)
		JSONR(c, 200, gin.H{"msg": "ok"})
	})
}
