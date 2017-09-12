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

	"github.com/gin-gonic/gin"
	"github.com/open-falcon/falcon-plus/modules/api/app/utils"
	"github.com/open-falcon/falcon-plus/modules/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	//session
	u := r.Group("/api/v1/user")
	u.GET("/auth_session", AuthSession)
	u.POST("/login", Login)
	u.GET("/logout", Logout)

	//user modify
	u.POST("/create", CreateUser)
	authapi := r.Group("/api/v1/user")
	authapi.Use(utils.AuthSessionMidd)
	authapi.GET("/current", UserInfo)
	authapi.GET("/u/:uid", GetUser)
	authapi.GET("/name/:user_name", GetUserByName)
	authapi.PUT("/update", UpdateCurrentUser)
	authapi.PUT("/cgpasswd", ChangePassword)
	authapi.GET("/users", UserList)
	authapi.GET("/u/:uid/in_teams", IsUserInTeams)
	authapi.GET("/u/:uid/teams", GetUserTeams)
	adminapi := r.Group("/api/v1/admin")
	adminapi.Use(utils.AuthSessionMidd)
	adminapi.PUT("/change_user_role", ChangeRoleOfUser)
	adminapi.PUT("/change_user_passwd", AdminChangePassword)
	adminapi.PUT("/change_user_profile", AdminChangeUserProfile)
	adminapi.DELETE("/delete_user", AdminUserDelete)

	//team
	authapi_team := r.Group("/api/v1")
	authapi_team.Use(utils.AuthSessionMidd)
	authapi_team.GET("/team", Teams)
	authapi_team.GET("/team/t/:team_id", GetTeam)
	authapi_team.GET("/team/name/:team_name", GetTeamByName)
	authapi_team.POST("/team", CreateTeam)
	authapi_team.PUT("/team", UpdateTeam)
	authapi_team.DELETE("/team/:team_id", DeleteTeam)
}
