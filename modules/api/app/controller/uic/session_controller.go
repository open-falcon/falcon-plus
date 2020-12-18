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

package uic

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	log "github.com/sirupsen/logrus"
)

type APILoginInput struct {
	Name     string `json:"name"  form:"name" binding:"required"`
	Password string `json:"password"  form:"password" binding:"required"`
}

type APIAdminLoginInput struct {
	Name string `json:"name"  form:"name" binding:"required"`
}

func Login(c *gin.Context) {
	inputs := APILoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "name or password is blank")
		return
	}
	user := uic.User{}
	db.Uic.Where(uic.User{Name: inputs.Name}).Find(&user)
	switch {
	case user.ID == 0:
		h.JSONR(c, badstatus, "no such user")
		return
	case user.Passwd != utils.HashIt(inputs.Password):
		h.JSONR(c, badstatus, "password error")
		return
	}
	var session uic.Session
	// response := map[string]string{}
	s := db.Uic.Table("session").Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, badstatus, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = user.ID
		db.Uic.Create(&session)
	}
	log.Debugf("session: %v", session)
	resp := struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{session.Sig, user.Name, user.IsAdmin()}
	h.JSONR(c, resp)
	return
}

func AdminLogin(c *gin.Context) {
	inputs := APIAdminLoginInput{}
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, "name is blank")
		return
	}
	name := inputs.Name

	user := uic.User{
		Name: name,
	}
	adminuser, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}

	db.Uic.Where(&user).Find(&user)
	switch {
	case user.ID == 0:
		h.JSONR(c, badstatus, "no such user")
		return
	case user.Role >= adminuser.Role:
		h.JSONR(c, badstatus, "API_USER not admin, no permissions can do this")
		return
	}
	var session uic.Session
	// response := map[string]string{}
	s := db.Uic.Table("session").Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, badstatus, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = user.ID
		db.Uic.Create(&session)
	}
	log.Debugf("session: %v", session)
	resp := struct {
		Sig   string `json:"sig,omitempty"`
		Name  string `json:"name,omitempty"`
		Admin bool   `json:"admin"`
	}{session.Sig, user.Name, user.IsAdmin()}
	h.JSONR(c, resp)
	return
}

func Logout(c *gin.Context) {
	wsession, err := h.GetSession(c)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	var session uic.Session
	var user uic.User
	db.Uic.Table("user").Where(uic.User{Name: wsession.Name}).Scan(&user)
	db.Uic.Table("session").Where("sig = ? AND uid = ?", wsession.Sig, user.ID).Scan(&session)

	if session.ID == 0 {
		h.JSONR(c, badstatus, "not found this kind of session in database.")
		return
	} else {
		r := db.Uic.Table("session").Delete(&session)
		if r.Error != nil {
			h.JSONR(c, badstatus, r.Error)
		}
		h.JSONR(c, "logout successful")
	}
	return
}

func AuthSession(c *gin.Context) {
	auth, err := h.SessionChecking(c)
	if err != nil || auth != true {
		h.JSONR(c, http.StatusUnauthorized, err)
		return
	}
	h.JSONR(c, "session is valid!")
	return
}

func CreateRoot(c *gin.Context) {
	password := c.DefaultQuery("password", "")
	if password == "" {
		h.JSONR(c, badstatus, "password is empty, please check it")
		return
	}
	password = utils.HashIt(password)
	user := uic.User{
		Name:   "root",
		Passwd: password,
	}
	dt := db.Uic.Table("user").Save(&user)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, "root created!")
	return
}
