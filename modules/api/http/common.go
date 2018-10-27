package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/file"
)

func configCommonRoutes() {
	routes.GET("/health", func(c *gin.Context) {
		c.Writer.Write([]byte("ok\n"))
	})

	routes.GET("/workdir", func(c *gin.Context) {
		c.Writer.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})
}
