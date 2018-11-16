// Copyright 2018 Xiaomi, Inc.
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

package alarm_manager

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    uint        `json:code`
	Data    interface{} `json:data`
	Message string      `json:message`
}

// Forwarder forwards the request to backend component and
// forwards the response from backend component to the caller of API.
func Forwarder(c *gin.Context) {
	request := c.Request
	if request == nil {
		c.JSON(http.StatusBadRequest, Resp{
			Code:    http.StatusBadRequest,
			Message: "request is nil",
		})
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	req, err := http.NewRequest(request.Method,
		server+request.URL.RequestURI(), bytes.NewBuffer(body))
	for k, v := range request.Header {
		for _, val := range v {
			req.Header.Add(k, val)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.Data(resp.StatusCode, "application/json", respBody)
	return
}
