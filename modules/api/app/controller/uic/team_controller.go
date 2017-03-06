package uic

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	h "github.com/open-falcon/falcon-plus/modules/api/app/helper"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

type CTeam struct {
	Team        uic.Team   `json:"team"`
	TeamCreator string     `json:"creator_name"`
	Users       []uic.User `json:"users"`
}

//support root as admin
func Teams(c *gin.Context) {
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
	query := c.DefaultQuery("q", ".+")
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	var dt *gorm.DB
	teams := []uic.Team{}
	if user.IsAdmin() {
		if limit != -1 && page != -1 {
			dt = db.Uic.Table("team").Raw(
				"select * from team where name regexp ? limit ?, ?", query, page, limit).Scan(&teams)
		} else {
			dt = db.Uic.Table("team").Where("name regexp ?", query).Scan(&teams)
		}
		err = dt.Error
	} else {
		//team creator and team member can manage the team
		dt = db.Uic.Raw(
			`select a.* from team as a, rel_team_user as b 
			where a.name regexp ? and a.id = b.tid and b.uid = ? 
			UNION select * from team where name regexp ? and creator = ?`,
			query, user.ID, query, user.ID).Scan(&teams)
		err = dt.Error
	}
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	outputs := []CTeam{}
	for _, t := range teams {
		cteam := CTeam{Team: t}
		user, err := t.Members()
		if err != nil {
			h.JSONR(c, badstatus, err)
			return
		}
		cteam.Users = user
		creatorName, err := t.GetCreatorName()
		if err != nil {
			log.Debug(err.Error())
		}
		cteam.TeamCreator = creatorName
		outputs = append(outputs, cteam)
	}
	h.JSONR(c, outputs)
	return
}

type APICreateTeamInput struct {
	Name    string  `json:"team_name" binding:"required"`
	Resume  string  `json:"resume"`
	UserIDs []int64 `json:"users"`
}

//every user can create a team
func CreateTeam(c *gin.Context) {
	var cteam APICreateTeamInput
	err := c.Bind(&cteam)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	team := uic.Team{
		Name:    cteam.Name,
		Resume:  cteam.Resume,
		Creator: user.ID,
	}
	dt := db.Uic.Table("team").Save(&team)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	var dt2 *gorm.DB
	if len(cteam.UserIDs) > 0 {
		for i := 0; i < len(cteam.UserIDs); i++ {
			dt2 = db.Uic.Save(&uic.RelTeamUser{Tid: team.ID, Uid: cteam.UserIDs[i]})
			if dt2.Error != nil {
				err = dt2.Error
				break
			}
		}
		if err != nil {
			h.JSONR(c, badstatus, err)
			return
		}
	}
	h.JSONR(c, fmt.Sprintf("team created! Afftect row: %d, Affect refer: %d", dt.RowsAffected, len(cteam.UserIDs)))
	return
}

type APIUpdateTeamInput struct {
	ID      int    `json:"team_id" binding:"required"`
	Resume  string `json:"resume"`
	UserIDs []int  `json:"users"`
}

// admin, team creator, team member can mangage the team
func UpdateTeam(c *gin.Context) {
	var cteam APIUpdateTeamInput
	err := c.Bind(&cteam)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}

	dt := db.Uic
	if user.IsAdmin() {
		dt = dt.Table("team").Where("id = ?", cteam.ID)
	} else {
		dt = dt.Raw(
			`select a.* from team as a, rel_team_user as b 
			where a.id = b.tid AND a.id = ? AND b.uid = ? 
			UNION select * from team where creator = ? AND id = ?`,
			cteam.ID, user.ID, user.ID, cteam.ID)
	}
	var team uic.Team
	dt = dt.Find(&team)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}

	err = bindUsers(db, cteam.ID, cteam.UserIDs)
	if err != nil {
		h.JSONR(c, badstatus, err)
	} else {
		h.JSONR(c, "team updated!")
	}
}

func bindUsers(db config.DBPool, tid int, users []int) (err error) {
	var dt *gorm.DB
	uids, err := utils.ArrIntToString(users)
	if err != nil {
		return
	}

	//delete unbind users
	var needDeleteMan []uic.RelTeamUser
	qPared := fmt.Sprintf("tid = %d AND NOT (uid IN (%v))", tid, uids)
	dt = db.Uic.Table("rel_team_user").Where(qPared).Find(&needDeleteMan)
	if dt.Error != nil {
		err = dt.Error
		return
	}
	if len(needDeleteMan) != 0 {
		for _, man := range needDeleteMan {
			dt = db.Uic.Delete(&man)
			if dt.Error != nil {
				err = dt.Error
				return
			}
		}
	}

	//insert bind users
	for _, i := range users {
		ur := uic.RelTeamUser{Tid: int64(tid), Uid: int64(i)}
		db.Uic.Table("rel_team_user").Where(&ur).Find(&ur)
		if ur.ID == 0 {
			dt = db.Uic.Table("rel_team_user").Save(&ur)
		} else {
			//if record exist, do next
			continue
		}
		if dt.Error != nil {
			err = dt.Error
			return
		}
	}
	return
}

type APIDeleteTeamInput struct {
	ID int64 `json:"team_id" binding:"required"`
}

//only admin or team creator can delete a team
func DeleteTeam(c *gin.Context) {
	var err error
	teamIdStr := c.Params.ByName("team_id")
	teamIdTmp, err := strconv.Atoi(teamIdStr)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	teamId := int64(teamIdTmp)
	if teamId == 0 {
		h.JSONR(c, badstatus, "team_id is empty")
		return
	} else if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	user, err := h.GetUser(c)
	if err != nil {
		h.JSONR(c, badstatus, err.Error())
		return
	}
	dt := db.Uic.Table("team")
	if user.IsAdmin() {
		dt = dt.Delete(&uic.Team{ID: teamId})
		err = dt.Error
	} else {
		team := uic.Team{
			ID:      teamId,
			Creator: user.ID,
		}
		dt = dt.Where(&team).Find(&team)
		if team.ID == 0 {
			err = errors.New("You don't have permission")
		} else if dt.Error != nil {
			err = dt.Error
		} else {
			db.Uic.Where("id = ?", teamId).Delete(&uic.Team{ID: teamId})
		}
	}
	var dt2 *gorm.DB
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, err)
		return
	} else {
		dt2 = db.Uic.Where("tid = ?", teamId).Delete(uic.RelTeamUser{})
	}
	h.JSONR(c, fmt.Sprintf("team %v is deleted. Affect row: %d / refer delete: %d", teamId, dt.RowsAffected, dt2.RowsAffected))
	return
}

type APIGetTeamOutput struct {
	uic.Team
	Users []uic.User `json:"users"`
}

func GetTeam(c *gin.Context) {
	team_id_str := c.Params.ByName("team_id")
	team_id, err := strconv.Atoi(team_id_str)
	if team_id == 0 {
		h.JSONR(c, badstatus, "team_id is empty")
		return
	} else if err != nil {
		h.JSONR(c, badstatus, err)
		return
	}
	team := uic.Team{ID: int64(team_id)}
	dt := db.Uic.Find(&team)
	if dt.Error != nil {
		h.JSONR(c, badstatus, dt.Error)
		return
	}
	var uidarr []uic.RelTeamUser
	db.Uic.Table("rel_team_user").Select("uid").Where(&uic.RelTeamUser{Tid: int64(team_id)}).Find(&uidarr)
	if err != nil {
		log.Debug(err.Error())
	}
	var resp APIGetTeamOutput
	resp.Team = team
	resp.Users = []uic.User{}
	if len(uidarr) != 0 {
		uids := ""
		for indx, v := range uidarr {
			if indx == 0 {
				uids = fmt.Sprintf("%v", v.Uid)
			} else {
				uids = fmt.Sprintf("%v,%v", uids, v.Uid)
			}
		}
		log.Debugf("uids:%s", uids)
		var users []uic.User
		db.Uic.Table("user").Where(fmt.Sprintf("id IN (%s)", uids)).Find(&users)
		resp.Users = users
	}
	h.JSONR(c, resp)
	return
}
