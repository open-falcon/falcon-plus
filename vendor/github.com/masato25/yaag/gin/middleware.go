package gin

import (
	"net/http/httptest"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/masato25/yaag/middleware"
	"github.com/masato25/yaag/yaag"
	"github.com/masato25/yaag/yaag/models"
)

func Document() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !yaag.IsOn() {
			return
		}
		writer := httptest.NewRecorder()
		apiCall := models.ApiCall{}
		middleware.Before(&apiCall, c.Request)
		c.Next()
		if writer.Code != 404 {
			apiCall.MethodType = c.Request.Method
			apiCall.CurrentPath = strings.Split(c.Request.RequestURI, "?")[0]
			body := ""
			if val, ok := c.Get("body_doc"); ok {
				body = val.(string)
			}
			apiCall.ResponseBody = body
			apiCall.ResponseCode = c.Writer.Status()
			headers := map[string]string{}
			for k, v := range c.Writer.Header() {
				log.Println(k, v)
				headers[k] = strings.Join(v, " ")
			}
			apiCall.ResponseHeader = headers
			//a custom api for fix broken header value
			apiCall.RequestHeader["Apitoken"] = c.Request.Header.Get("Apitoken")
			go yaag.GenerateHtml(&apiCall)
		}
	}
}
