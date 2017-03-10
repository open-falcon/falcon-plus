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
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/httplib"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
)

const (
	MISSO_NAME = "misso"
)

type missoAuth struct {
	RedirectURL string

	CookieSecretKey string
	missoAuthDomain string
	BrokerName      string
	SecretKey       string
	Credential      string
}

func init() {
	models.RegisterAuth(MISSO_NAME, &missoAuth{})
}

func (p *missoAuth) Init(conf *falcon.ConfCtrl) error {
	p.RedirectURL = conf.Ctrl.Str(falcon.C_MISSO_REDIRECT_URL)
	p.CookieSecretKey = "secret-key-for-encrypt-cookie"
	p.missoAuthDomain = "http://sso.pt.xiaomi.com"
	p.BrokerName = "test"
	p.SecretKey = "test"
	return nil
}

func (p *missoAuth) Verify(c interface{}) (bool, string, error) {
	return false, "", models.EPERM
}

func (p *missoAuth) AuthorizeUrl(c interface{}) string {
	ctx := c.(*context.Context)

	p.Credential, _ = p.GenerateCredential()
	ctx.SetCookie("broker_cookie", p.Credential)

	v := url.Values{}
	v.Set("callback", p.RedirectURL)

	url, err := p.GetLoginUrl()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s&%s", url, v.Encode())
}

func (p *missoAuth) CallBack(c interface{}) (uuid string, err error) {
	ctx := c.(*context.Context)

	remote_ip := models.GetIPAdress(ctx.Input.Context.Request)

	user_name, result := ctx.GetSecureCookie(p.CookieSecretKey, "user_name")
	broker_cookie := ctx.GetCookie("broker_cookie")

	//If can get user_name from cookie, user have logined
	if result == true {
		uuid = fmt.Sprintf("%s@%s", user_name, MISSO_NAME)
		return
	} else {
		if broker_cookie == "" {
			//cannot get broker_cookie, may first open in browser, or use service_account
			//try to get user_name from sso, may be login use service account
			authorization := ctx.Input.Header("Authorization")
			if authorization != "" {
				uuid, err = p.GetServiceUser(authorization, remote_ip)
				if err == nil && user_name != "" {
					uuid = fmt.Sprintf("%s@%s", uuid, MISSO_NAME)
					return
				}
			}
		} else {
			p.Credential = broker_cookie
			if user_name, _ = p.GetUser(); user_name != "" {
				uuid = fmt.Sprintf("%s@%s", user_name, MISSO_NAME)
				ctx.SetSecureCookie(p.CookieSecretKey,
					"user_name", user_name)
				return
			}
		}
	}
	err = models.ErrLogin
	return
}

/***********************
 * from sso_client.go
 ***********************/
func (p *missoAuth) GetUser() (string, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/user",
		p.missoAuthDomain, p.BrokerName, p.Credential)
	resp, err := httplib.Get(url).String()
	var resp_js map[string]string
	err = json.Unmarshal([]byte(resp), &resp_js)
	return resp_js["user_name"], err
}

func (p *missoAuth) GetServiceUser(authorization, user_ip string) (string,
	error) {

	auth_len := strings.Split(authorization, ";")
	if len(auth_len) != 3 {
		return "", fmt.Errorf("authorization wrong")
	}
	url := fmt.Sprintf("%s/mias/api/user_ip/%s/auth/%s/username",
		p.missoAuthDomain, user_ip, authorization)
	resp, err := httplib.Get(url).String()
	if err != nil {
		return "", err
	}
	var resp_js map[string]string
	err = json.Unmarshal([]byte(resp), &resp_js)
	if err != nil {
		return "", fmt.Errorf(resp)
	}

	return resp_js["user_name"], nil
}

func (p *missoAuth) IsLogin() (bool, error) {
	url := fmt.Sprintf("%s/login/broker/%s/broker_cookies/%s/check",
		p.missoAuthDomain, p.BrokerName, p.Credential)
	resp, err := httplib.Get(url).String()
	if err != nil {
		return false, err
	}

	return resp == "1", nil
}

func (p *missoAuth) GetLogoutUrl() string {
	return fmt.Sprintf("%s/login/logout?broker_name=%s",
		p.missoAuthDomain, p.BrokerName)
}

func (p *missoAuth) GetLoginUrl() (string, error) {
	if p.Credential == "" {
		_, err := p.GenerateCredential()
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s/login?broker_cookies=%s",
		p.missoAuthDomain, p.Credential), nil
}

func (p *missoAuth) GenerateCredential() (string, error) {
	url := p.missoAuthDomain + "/login/broker_cookies"
	req := httplib.Post(url)
	req.Param("broker_name", p.BrokerName)
	req.Param("secret_key", p.SecretKey)
	resp, err := req.SetTimeout(3*time.Second,
		3*time.Second).String()
	p.Credential = resp
	return resp, err
}
