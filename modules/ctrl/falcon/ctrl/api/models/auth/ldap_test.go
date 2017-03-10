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

import "testing"

func TestVerify(t *testing.T) {

	addr := "localhost:389"
	baseDN := "dc=yubo,dc=org"
	username := "yubo"
	password := "12341234"
	bindusername := "cn=admin,dc=yubo,dc=org"
	bindpassword := "12341234"
	filter := "(&(objectClass=posixAccount)(cn=%s))"

	success, userDN, err := ldapUserAuthentication(addr, baseDN, filter, username, password, bindusername, bindpassword, false)
	t.Log("success:", success)
	t.Log("userDN:", userDN)
	t.Log("err:", err)
}
