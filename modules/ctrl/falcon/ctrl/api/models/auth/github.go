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
package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego/context"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	GITHUB_BASE_URL = "https://api.github.com"
	GITHUB_NAME     = "github"
)

type githubAuth struct {
	config oauth2.Config
}

func init() {
	models.RegisterAuth(GITHUB_NAME, &githubAuth{})
}

func (p *githubAuth) Init(conf *falcon.ConfCtrl) error {
	p.config = oauth2.Config{
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
		ClientID:     conf.Ctrl.Str(falcon.C_GITHUB_CLIENT_ID),
		ClientSecret: conf.Ctrl.Str(falcon.C_GITHUB_CLIENT_SECRET),
		RedirectURL:  conf.Ctrl.Str(falcon.C_GITHUB_REDIRECT_URL),
	}
	return nil
}

func (p *githubAuth) Verify(c interface{}) (bool, string, error) {
	return false, "", models.EPERM
}

func (p *githubAuth) AuthorizeUrl(c interface{}) string {
	ctx := c.(*context.Context)

	v := url.Values{}
	v.Set("cb", ctx.Input.Query("cb"))

	conf := p.config
	conf.RedirectURL = fmt.Sprintf("%s?%s", conf.RedirectURL, v.Encode())
	return conf.AuthCodeURL(models.RandString(8))
}

func (p *githubAuth) CallBack(c interface{}) (uuid string, err error) {
	var user githubUser

	r := c.(*context.Context).Request
	q := r.URL.Query()

	if errType := q.Get("error"); errType != "" {
		err = fmt.Errorf("%s:%s", errType, q.Get("error_description"))
		return
	}

	ctx := r.Context()

	token, err := p.config.Exchange(ctx, q.Get("code"))
	if err != nil {
		err = fmt.Errorf("github: failed to get token: %v", err)
		return
	}

	client := p.config.Client(ctx, token)

	user, err = p.user(r, client)
	if err != nil {
		err = fmt.Errorf("github: get user: %v", err)
		return
	}

	//beego.Debug(fmt.Sprintf("name:%s login:%s id:%d email:%s", user.Name, user.Login, user.ID, user.Email))
	uuid = fmt.Sprintf("%s@%s", user.Login, GITHUB_NAME)
	return
}

type githubUser struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// user queries the GitHub API for profile information using the provided client. The HTTP
// client is expected to be constructed by the golang.org/x/oauth2 package, which inserts
// a bearer token as part of the request.
func (c *githubAuth) user(r *http.Request, client *http.Client) (githubUser, error) {
	var u githubUser
	req, err := http.NewRequest("GET", GITHUB_BASE_URL+"/user", nil)
	if err != nil {
		return u, fmt.Errorf("github: new req: %v", err)
	}
	req = req.WithContext(r.Context())
	resp, err := client.Do(req)
	if err != nil {
		return u, fmt.Errorf("github: get URL %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return u, fmt.Errorf("github: read body: %v", err)
		}
		return u, fmt.Errorf("%s: %s", resp.Status, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode response: %v", err)
	}
	return u, nil
}
