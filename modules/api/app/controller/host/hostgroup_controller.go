package host

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	f "github.com/open-falcon/falcon-plus/modules/api/app/model/falcon_portal"
	u "github.com/open-falcon/falcon-plus/modules/api/app/utils"
)

func GetHostGroups(c *gin.Context) {
	var (
		limit int
		page  int
		err   error
	)
	pageTmp := c.DefaultQuery("page", "")
	limitTmp := c.DefaultQuery("limit", "")
	q := c.DefaultQuery("q", ".+")
	page, limit, err = h.PageParser(pageTmp, limitTmp)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	var hostgroups []f.HostGroup
	var dt *gorm.DB
	if limit != -1 && page != -1 {
		dt = db.Falcon.Raw(fmt.Sprintf("SELECT * from grp  where grp_name regexp '%s' limit %d,%d", q, page, limit)).Scan(&hostgroups)
	} else {
		dt = db.Falcon.Table("grp").Where("grp_name regexp ?", q).Find(&hostgroups)
	}
	if dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, hostgroups)
	return
}

type APICrateHostGroup struct {
	Name string `json:"name" binding:"required"`
}

func CrateHostGroup(c *gin.Context) {
	var inputs APICrateHostGroup
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	hostgroup := f.HostGroup{Name: inputs.Name, CreateUser: user.Name, ComeFrom: 1}
	if dt := db.Falcon.Create(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, hostgroup)
	return
}

type APIBindHostToHostGroupInput struct {
	Hosts       []string `json:"hosts" binding:"required"`
	HostGroupID int64    `json:"hostgroup_id" binding:"required"`
}

func BindHostToHostGroup(c *gin.Context) {
	var inputs APIBindHostToHostGroupInput
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	hostgroup := f.HostGroup{ID: inputs.HostGroupID}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	if !user.IsAdmin() && hostgroup.CreateUser != user.Name {
		h.JSONR(c, expecstatus, "You don't have permission.")
		return
	}
	tx := db.Falcon.Begin()
	if dt := tx.Where("grp_id = ?", hostgroup.ID).Delete(&f.GrpHost{}); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("delete grp_host got error: %v", dt.Error))
		dt.Rollback()
		return
	}
	var ids []int64
	for _, host := range inputs.Hosts {
		ahost := f.Host{Hostname: host}
		var id int64
		var ok bool
		if id, ok = ahost.Existing(); ok {
			ids = append(ids, id)
		} else {
			if dt := tx.Save(&ahost); dt.Error != nil {
				h.JSONR(c, expecstatus, dt.Error)
				tx.Rollback()
				return
			}
			id = ahost.ID
			ids = append(ids, id)
		}
		if dt := tx.Debug().Create(&f.GrpHost{GrpID: hostgroup.ID, HostID: id}); dt.Error != nil {
			h.JSONR(c, expecstatus, fmt.Sprintf("create grphost got error: %s , grp_id: %v, host_id: %v", dt.Error, hostgroup.ID, id))
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("%v bind to hostgroup: %v", ids, hostgroup.ID))
	return
}

type APIUnBindAHostToHostGroup struct {
	HostID      int64 `json:"host_id" binding:"required"`
	HostGroupID int64 `json:"hostgroup_id" binding:"required"`
}

func UnBindAHostToHostGroup(c *gin.Context) {
	var inputs APIUnBindAHostToHostGroup
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	hostgroup := f.HostGroup{ID: inputs.HostGroupID}
	if !user.IsAdmin() {
		if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		}
		if hostgroup.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			return
		}
	}
	if dt := db.Falcon.Where("grp_id = ? AND host_id = ?", inputs.HostGroupID, inputs.HostID).Delete(&f.GrpHost{}); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("unbind host:%v of hostgroup: %v", inputs.HostID, inputs.HostGroupID))
	return
}

func DeleteHostGroup(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if !user.IsAdmin() {
		if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
			h.JSONR(c, badstatus, dt.Error)
			return
		}
		if hostgroup.CreateUser != user.Name {
			h.JSONR(c, badstatus, "You don't have permission!")
			return
		}
	}
	tx := db.Falcon.Begin()
	//delete hostgroup referance of grp_host table
	if dt := tx.Where("grp_id = ?", grpID).Delete(&f.GrpHost{}); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("delete grp_host got error: %v", dt.Error))
		dt.Rollback()
		return
	}
	//delete plugins of hostgroup
	if dt := tx.Where("grp_id = ?", grpID).Delete(&f.Plugin{}); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("delete plugins got error: %v", dt.Error))
		dt.Rollback()
		return
	}
	//delete aggreators of hostgroup
	if dt := tx.Where("grp_id = ?", grpID).Delete(&f.Cluster{}); dt.Error != nil {
		h.JSONR(c, expecstatus, fmt.Sprintf("delete aggreators got error: %v", dt.Error))
		dt.Rollback()
		return
	}
	//finally delete hostgroup
	if dt := tx.Delete(&f.HostGroup{ID: int64(grpID)}); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		tx.Rollback()
		return
	}
	tx.Commit()
	h.JSONR(c, fmt.Sprintf("hostgroup:%v has been deleted", grpID))
	return
}

func GetHostGroup(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	q := c.DefaultQuery("q", ".+")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	hosts := []f.Host{}
	grpHosts := []f.GrpHost{}
	if dt := db.Falcon.Where("grp_id = ?", grpID).Find(&grpHosts); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	for _, grph := range grpHosts {
		var host f.Host
		db.Falcon.Find(&host, grph.HostID)
		if host.ID != 0 {
			if ok, err := regexp.MatchString(q, host.Hostname); ok == true && err == nil {
				hosts = append(hosts, host)
			}
		}
	}
	h.JSONR(c, map[string]interface{}{
		"hostgroup": hostgroup,
		"hosts":     hosts,
	})
	return
}

type APIHostGroupInputs struct {
	Name string `json:"grp_name" binding:"required"`
	//create_user string `json:"create_user" binding:"required"`
}
func PutHostGroup(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	var inputs APIHostGroupInputs
	err = c.BindJSON(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, badstatus, err)
		return
	case u.HasDangerousCharacters(inputs.Name):
		h.JSONR(c, badstatus, "grp_name is invalid")
		return
	}
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	hostgroup.Name = inputs.Name
	uhostgroup := map[string]interface{}{
		"grp_name":     hostgroup.Name,
		"create_user":  hostgroup.CreateUser,
		"come_from":    hostgroup.ComeFrom,
	}
	dt := db.Falcon.Model(&hostgroup).Where("id = ?", grpID).Update(uhostgroup)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, "hostgroup profile updated")
	return
}

type APIBindTemplateToGroupInputs struct {
	TplID int64 `json:"tpl_id"`
	GrpID int64 `json:"grp_id"`
}

func BindTemplateToGroup(c *gin.Context) {
	var inputs APIBindTemplateToGroupInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	grpTpl := f.GrpTpl{
		GrpID: inputs.GrpID,
		TplID: inputs.TplID,
	}
	db.Falcon.Where("grp_id = ? and tpl_id = ?", inputs.GrpID, inputs.TplID).Find(&grpTpl)
	if grpTpl.BindUser != "" {
		h.JSONR(c, badstatus, errors.New("this binding already existing, reject!"))
		return
	}
	grpTpl.BindUser = user.Name
	if dt := db.Falcon.Save(&grpTpl); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, grpTpl)
	return
}

type APIUnBindTemplateToGroupInputs struct {
	TplID int64 `json:"tpl_id"`
	GrpID int64 `json:"grp_id"`
}

func UnBindTemplateToGroup(c *gin.Context) {
	var inputs APIUnBindTemplateToGroupInputs
	if err := c.Bind(&inputs); err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, _ := h.GetUser(c)
	grpTpl := f.GrpTpl{
		GrpID: inputs.GrpID,
		TplID: inputs.TplID,
	}
	db.Falcon.Where("grp_id = ? and tpl_id = ?", inputs.GrpID, inputs.TplID).Find(&grpTpl)
	switch {
	case !user.IsAdmin() && grpTpl.BindUser != user.Name:
		h.JSONR(c, badstatus, errors.New("You don't have permission can do this."))
		return
	}
	if dt := db.Falcon.Where("grp_id = ? and tpl_id = ?", inputs.GrpID, inputs.TplID).Delete(&grpTpl); dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	h.JSONR(c, fmt.Sprintf("template: %v is unbind of HostGroup: %v", inputs.TplID, inputs.GrpID))
	return
}

func GetTemplateOfHostGroup(c *gin.Context) {
	grpIDtmp := c.Params.ByName("host_group")
	if grpIDtmp == "" {
		h.JSONR(c, badstatus, "grp id is missing")
		return
	}
	grpID, err := strconv.Atoi(grpIDtmp)
	if err != nil {
		log.Debugf("grpIDtmp: %v", grpIDtmp)
		h.JSONR(c, badstatus, err)
		return
	}
	hostgroup := f.HostGroup{ID: int64(grpID)}
	if dt := db.Falcon.Find(&hostgroup); dt.Error != nil {
		h.JSONR(c, expecstatus, dt.Error)
		return
	}
	grpTpls := []f.GrpTpl{}
	Tpls := []f.Template{}
	db.Falcon.Where("grp_id = ?", grpID).Find(&grpTpls)
	if len(grpTpls) != 0 {
		tips := []int64{}
		for _, t := range grpTpls {
			tips = append(tips, t.TplID)
		}
		tipsStr, _ := u.ArrInt64ToString(tips)
		db.Falcon.Where(fmt.Sprintf("id in (%s)", tipsStr)).Find(&Tpls)
	}
	h.JSONR(c, map[string]interface{}{
		"hostgroup": hostgroup,
		"templates": Tpls,
	})
	return
}
