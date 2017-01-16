package uic

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	h "github.com/masato25/owl_backend/app/helper"
	"github.com/masato25/owl_backend/app/model/uic"
	"github.com/masato25/owl_backend/app/utils"
)

func Login(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	password := c.DefaultQuery("password", "")

	if name == "" || password == "" {
		h.JSONR(c, badstatus, "name or password is blank")
		return
	}
	user := uic.User{
		Name: name,
	}
	db.Uic.Where(&user).Find(&user)
	switch {
	case user.Name == "":
		h.JSONR(c, badstatus, "no such user")
		return
	case user.Passwd != utils.HashIt(password):
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
	h.JSONR(c, "session is vaild!")
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
