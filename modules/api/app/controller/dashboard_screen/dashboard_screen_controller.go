package dashboard_screen

import (
	"fmt"
	"github.com/gin-gonic/gin"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	m "github.com/open-falcon/falcon-plus/modules/api/app/model/dashboard"
	"strconv"
)

func ScreenCreate(c *gin.Context) {
	pid := c.DefaultPostForm("pid", "0")
	name := c.DefaultPostForm("name", "")
	if name == "" {
		h.JSONR(c, badstatus, "empty name")
		return
	}

	ipid, err := strconv.Atoi(pid)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen pid")
		return
	}

	dt := db.Dashboard.Exec("insert ignore into dashboard_screen (pid, name) values(?, ?)", ipid, name)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	var lid []int
	dt = db.Dashboard.Table("dashboard_screen").Select("id").Where("pid = ? and name = ?", ipid, name).Limit(1).Pluck("id", &lid)
	if dt.Error != nil || len(lid) == 0 {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	if len(lid) == 0 {
		h.JSONR(c, badstatus, fmt.Sprintf("no such screen where name=%s", name))
		return
	}
	sid := lid[0]

	h.JSONR(c, map[string]interface{}{"pid": ipid, "id": sid, "name": name})
}

func ScreenGet(c *gin.Context) {
	id := c.Param("screen_id")

	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}

	screen := m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Where("id = ?", sid).First(&screen)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, screen)
}

func ScreenGetsByPid(c *gin.Context) {
	id := c.Param("pid")

	pid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen pid")
		return
	}

	screens := []m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Where("pid = ?", pid).Find(&screens)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, screens)
}

func ScreenGetsAll(c *gin.Context) {
	limit := c.DefaultQuery("limit", "500")
	screens := []m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Limit(limit).Find(&screens)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, screens)
}

func ScreenDelete(c *gin.Context) {
	id := c.Param("screen_id")

	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}

	screen := m.DashboardScreen{}
	dt := db.Dashboard.Table("dashboard_screen").Where("id = ?", sid).Delete(&screen)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, "ok")
}

func ScreenUpdate(c *gin.Context) {
	id := c.Param("screen_id")

	sid, err := strconv.Atoi(id)
	if err != nil {
		h.JSONR(c, badstatus, "invalid screen id")
		return
	}

	new_data := map[string]interface{}{}
	pid := c.PostForm("pid")
	name := c.PostForm("name")
	if name != "" {
		new_data["name"] = name
	}

	if pid != "" {
		ipid, err := strconv.Atoi(pid)
		if err != nil {
			h.JSONR(c, badstatus, "invalid screen pid")
			return
		}
		new_data["pid"] = ipid
	}

	dt := db.Dashboard.Table("dashboard_screen").Where("id = ?", sid).Update(new_data)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	h.JSONR(c, "ok")
}
