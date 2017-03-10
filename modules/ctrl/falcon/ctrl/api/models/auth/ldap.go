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
	"crypto/tls"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/controllers"
	"github.com/open-falcon/falcon-plus/modules/ctrl/falcon/ctrl/api/models"
	"gopkg.in/ldap.v2"
)

type ldapAuth struct {
	addr    string
	baseDN  string
	bindDN  string
	bindPwd string
	filter  string
	tls     bool
}

const (
	LDAP_NAME = "ldap"
)

func init() {
	models.RegisterAuth(LDAP_NAME, &ldapAuth{})
}

func (p *ldapAuth) Init(conf *falcon.ConfCtrl) error {
	p.addr = conf.Ctrl.Str(falcon.C_LDAP_ADDR)
	p.baseDN = conf.Ctrl.Str(falcon.C_LDAP_BASE_DN)
	p.bindDN = conf.Ctrl.Str(falcon.C_LDAP_BIND_DN)
	p.bindPwd = conf.Ctrl.Str(falcon.C_LDAP_BIND_PWD)
	p.filter = conf.Ctrl.Str(falcon.C_LDAP_FILTER)
	return nil
}

func (p *ldapAuth) Verify(_c interface{}) (bool, string, error) {
	c := _c.(*controllers.AuthController)
	username := c.GetString("username")
	password := c.GetString("password")

	if beego.BConfig.RunMode == "dev" && username == "test" {
		return true, "test", nil
	}

	success, uuid, err := ldapUserAuthentication(p.addr, p.baseDN, p.filter,
		username, password,
		p.bindDN, p.bindPwd, p.tls)
	if success {
		uuid = fmt.Sprintf("%s@%s", uuid, LDAP_NAME)
	}

	return success, uuid, err
}

func (p *ldapAuth) AuthorizeUrl(c interface{}) string {
	return ""
}

func (p *ldapAuth) CallBack(c interface{}) (uuid string, err error) {
	return "", models.EPERM
}

func ldapUserAuthentication(addr, baseDN, filter, username, password, bindusername, bindpassword string, TLS bool) (success bool, userDN string, err error) {
	var (
		sr *ldap.SearchResult
	)

	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()

	// Reconnect with TLS
	if TLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			beego.Warning(err)
		}
	}

	// First bind with a read only user
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		return
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(filter, username),
		[]string{"dn"},
		nil,
	)

	sr, err = l.Search(searchRequest)
	if err != nil {
		return
	}

	if len(sr.Entries) != 1 {
		return
	}

	userDN = sr.Entries[0].DN

	// Bind as the user to verify their password
	err = l.Bind(userDN, password)
	if err != nil {
		return
	}
	return true, userDN, err

}
