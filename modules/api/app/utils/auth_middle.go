// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/spf13/viper"
)

func AuthSessionMidd(c *gin.Context) {
	auth, err := h.SessionChecking(c)
	if !viper.GetBool("skip_auth") {
		if err != nil || auth != true {
			log.Debugf("error: %v, auth: %v", err.Error(), auth)
			c.Set("auth", auth)
			h.JSONR(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}
	}
	c.Set("auth", auth)
}

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Add("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Apitoken")
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(200)
		} else {
			context.Next()
		}
	}
}
