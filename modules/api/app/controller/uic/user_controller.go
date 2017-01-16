package uic

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/masato25/owl_backend/app/helper"
	"github.com/masato25/owl_backend/app/model/uic"
	"github.com/masato25/owl_backend/app/utils"
)

type APIUserInput struct {
	Name   string `json:"name" binding:"required"`
	Cnname string `json:"cnname" binding:"required"`
	Passwd string `json:"password" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Phone  string `json:"phone"`
	IM     string `json:"im"`
	QQ     string `json:"qq"`
}

func CreateUser(c *gin.Context) {
	var inputs APIUserInput
	err := c.Bind(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, http.StatusBadRequest, err)
		return
	case utils.HasDangerousCharacters(inputs.Cnname):
		h.JSONR(c, http.StatusBadRequest, "name pattern is invalid")
		return
	}
	var user uic.User
	db.Uic.Table("user").Where("name = ?", inputs.Name).Scan(&user)
	if user.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, "name is already existing")
		return
	}
	password := utils.HashIt(inputs.Passwd)
	user = uic.User{
		Name:   inputs.Name,
		Passwd: password,
		Cnname: inputs.Cnname,
		Email:  inputs.Email,
		Phone:  inputs.Phone,
		IM:     inputs.IM,
		QQ:     inputs.QQ,
	}

	dt := db.Uic.Table("user").Create(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, dt.Error)
		return
	}

	var session uic.Session
	response := map[string]string{}
	s := db.Uic.Table("session").Where("uid = ?", user.ID).Scan(&session)
	if s.Error != nil && s.Error.Error() != "record not found" {
		h.JSONR(c, http.StatusBadRequest, s.Error)
		return
	} else if session.ID == 0 {
		session.Sig = utils.GenerateUUID()
		session.Expired = int(time.Now().Unix()) + 3600*24*30
		session.Uid = user.ID
		db.Uic.Create(&session)
	}
	log.Debugf("%v", session)
	response["sig"] = session.Sig
	response["name"] = user.Name
	h.JSONR(c, http.StatusOK, response)
	return
}

type APIUserUpdateInput struct {
	Cnname string `json:"cnname" binding:"required"`
	Email  string `json:"email" binding:"required"`
	Phone  string `json:"phone"`
	IM     string `json:"im"`
	QQ     string `json:"qq"`
}

func UpdateUser(c *gin.Context) {
	var inputs APIUserUpdateInput
	err := c.BindJSON(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, http.StatusExpectationFailed, err)
		return
	case utils.HasDangerousCharacters(inputs.Cnname):
		h.JSONR(c, http.StatusBadRequest, "name pattern is invalid")
		return
	}
	websession, _ := h.GetSession(c)
	user := uic.User{}
	db.Uic.Table("user").Where("name = ?", websession.Name).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "name is not existing")
		return
	}
	uid := user.ID
	uuser := map[string]interface{}{
		"Cnname": inputs.Cnname,
		"Email":  inputs.Email,
		"Phone":  inputs.Phone,
		"IM":     inputs.IM,
		"QQ":     inputs.QQ,
	}
	dt := db.Uic.Model(&user).Where("id = ?", uid).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, "user info updated")
	return
}

type APICgPassedInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	var inputs APICgPassedInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
	}
	websession, _ := h.GetSession(c)
	user := uic.User{Name: websession.Name}

	dt := db.Uic.Where(&user).Find(&user)
	switch {
	case dt.Error != nil:
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	case user.Passwd != utils.HashIt(inputs.OldPassword):
		h.JSONR(c, http.StatusBadRequest, "oldPassword is not match current one")
		return
	case user.Passwd != "":
		h.JSONR(c, http.StatusBadRequest, "Password can not be blank")
		return
	}

	user.Passwd = utils.HashIt(inputs.NewPassword)
	dt = db.Uic.Save(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, http.StatusOK, "password updated!")
	return
}

func UserInfo(c *gin.Context) {
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, err)
		return
	}
	h.JSONR(c, http.StatusOK, user)
	return
}

func GetUser(c *gin.Context) {
	uidtmp := c.Params.ByName("uid")
	if uidtmp == "" {
		h.JSONR(c, badstatus, "user id is missing")
		return
	}
	uid, err := strconv.Atoi(uidtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, err := h.GetUser(c)
	if !user.IsAdmin() {
		h.JSONR(c, badstatus, "only admin user can do this.")
		return
	}
	fuser := uic.User{ID: int64(uid)}
	if dt := db.Uic.Find(&fuser); dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, fuser)
	return
}

type APIAdminUserDeleteInput struct {
	UserID int `json:"user_id" binding:"required"`
}

//admin usage
func AdminUserDelete(c *gin.Context) {
	var inputs APIAdminUserDeleteInput
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	cuser, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, err)
		return
	} else if !cuser.IsAdmin() {
		h.JSONR(c, http.StatusBadRequest, "you don't have permission!")
		return
	}
	dt := db.Uic.Delete(&uic.User{}, inputs.UserID)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("user %v has been delete, affect row: %v", inputs.UserID, dt.RowsAffected))
	return
}

type APIAdminChangePassword struct {
	UserID int    `json:"user_id" binding:"required"`
	Passwd string `json:"password" binding:"required"`
}

//admin usage
func AdminChangePassword(c *gin.Context) {
	var inputs APIAdminChangePassword
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	websession, _ := h.GetSession(c)
	user := uic.User{Name: websession.Name}
	dt := db.Uic.Where(&user).Find(&user)
	switch {
	case dt.Error != nil:
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	case !user.IsAdmin():
		h.JSONR(c, http.StatusBadRequest, "you don't have permission!")
		return
	}

	user.Passwd = utils.HashIt(inputs.Passwd)
	dt = db.Uic.Save(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, http.StatusOK, "password updated!")
	return
}

func UserList(c *gin.Context) {
	// remove admin checking
	// websession, _ := h.GetSession(c)
	// user := uic.User{Name: websession.Name}
	// dt := db.Uic.Where(&user).Find(&user)
	// switch {
	// case dt.Error != nil:
	// 	h.JSONR(c, http.StatusExpectationFailed, dt.Error)
	// 	return
	// case !user.IsAdmin():
	// 	h.JSONR(c, http.StatusBadRequest, "you don't have permission!")
	// 	return
	// }
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	q := c.DefaultQuery("q", ".+")
	var user []uic.User
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Uic.Raw(
			fmt.Sprintf("select * from user where name regexp '%s' limit %d,%d", q, page, limit)).Scan(&user)
	} else {
		dt = db.Uic.Table("user").Where("name regexp ?", q).Scan(&user)
	}
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, user)
	return
}

//admin usage
type APIRoleUpdate struct {
	UserID int64  `json:"user_id" binding:"required"`
	Admin  string `json:"admin" binding:"required"`
}

func ChangeRuleOfUser(c *gin.Context) {
	var inputs APIRoleUpdate
	err := c.Bind(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusBadRequest, err)
		return
	}
	cuser, err := h.GetUser(c)
	switch {
	case err != nil:
		h.JSONR(c, http.StatusBadRequest, err)
		return
	case !cuser.IsAdmin():
		h.JSONR(c, http.StatusBadRequest, "you don't have permission!")
		return
	}
	var user uic.User
	db.Uic.Find(&user, inputs.UserID)
	switch inputs.Admin {
	case "yes":
		user.Role = 1
	case "no":
		user.Role = 0
	}
	log.Debugf("inputs got %v, user: %v", inputs, user)
	dt := db.Uic.Save(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("user role update sccuessful, affect row: %v", dt.RowsAffected))
	return
}
