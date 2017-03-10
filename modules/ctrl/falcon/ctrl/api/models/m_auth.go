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
package models

import "github.com/open-falcon/falcon-plus/modules/ctrl/falcon"

var (
	allAuths = make(map[string]AuthInterface)
	Auths    = make(map[string]AuthInterface)
)

type Auth struct {
	Method string
	Arg1   string
	Arg2   string
}

type AuthInterface interface {
	Init(conf *falcon.ConfCtrl) error
	Verify(c interface{}) (success bool, uuid string, err error)
	AuthorizeUrl(ctx interface{}) string
	CallBack(ctx interface{}) (uuid string, err error)
}

func RegisterAuth(name string, p AuthInterface) error {
	if _, ok := allAuths[name]; ok {
		return ErrExist
	} else {
		allAuths[name] = p
		return nil
	}
}
