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

package utils

import (
	"golang.org/x/crypto/bcrypt"

	log "github.com/Sirupsen/logrus"
)

func HashIt(passwd string) (hashed string) {
	if bs, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err != nil {
		log.Errorf("bcrypt fail, error:%v", err)
		hashed = ""
	} else {
		hashed = string(bs)
	}
	return
}
