package api

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/uic"
	"github.com/toolkits/container/set"
	"github.com/toolkits/net/httplib"
	"strings"
	"sync"
	"time"
)

type APIGetTeamOutput struct {
	uic.Team
	Users       []*uic.User `json:"users"`
	TeamCreator string      `json:"creator_name"`
}

type UsersCache struct {
	sync.RWMutex
	M map[string][]*uic.User
}

var Users = &UsersCache{M: make(map[string][]*uic.User)}

func (this *UsersCache) Get(team string) []*uic.User {
	this.RLock()
	defer this.RUnlock()
	val, exists := this.M[team]
	if !exists {
		return nil
	}

	return val
}

func (this *UsersCache) Set(team string, users []*uic.User) {
	this.Lock()
	defer this.Unlock()
	this.M[team] = users
}

func UsersOf(team string) []*uic.User {
	users := CurlUic(team)

	if users != nil {
		Users.Set(team, users)
	} else {
		users = Users.Get(team)
	}

	return users
}

func GetUsers(teams string) map[string]*uic.User {
	userMap := make(map[string]*uic.User)
	arr := strings.Split(teams, ",")
	for _, team := range arr {
		if team == "" {
			continue
		}

		users := UsersOf(team)
		if users == nil {
			continue
		}

		for _, user := range users {
			userMap[user.Name] = user
		}
	}
	return userMap
}

// return phones, emails
func ParseTeams(teams string) ([]string, []string) {
	if teams == "" {
		return []string{}, []string{}
	}

	userMap := GetUsers(teams)
	phoneSet := set.NewStringSet()
	mailSet := set.NewStringSet()
	for _, user := range userMap {
		phoneSet.Add(user.Phone)
		mailSet.Add(user.Email)
	}
	return phoneSet.ToSlice(), mailSet.ToSlice()
}

func CurlUic(team string) []*uic.User {
	if team == "" {
		return []*uic.User{}
	}

	uri := fmt.Sprintf("%s/api/v1/team/name/%s", g.Config().Api.PlusApi, team)
	req := httplib.Get(uri).SetTimeout(2*time.Second, 10*time.Second)
	token, _ := json.Marshal(map[string]string{
		"name": "falcon-alarm",
		"sig":  g.Config().Api.PlusApiToken,
	})
	req.Header("Apitoken", string(token))

	var team_users APIGetTeamOutput
	err := req.ToJson(&team_users)
	if err != nil {
		log.Errorf("curl %s fail: %v", uri, err)
		return nil
	}

	return team_users.Users
}
