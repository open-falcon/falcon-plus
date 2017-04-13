package http

import (
	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/toolkits/file"
)

func Version(c *gin.Context) {
	c.String(200, g.VERSION)
}

func Health(c *gin.Context) {
	c.String(200, "ok")
}

func Workdir(c *gin.Context) {
	c.String(200, file.SelfDir())
}
