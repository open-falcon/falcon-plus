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
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

func UniqSet(list []interface{}) []interface{} {
	mapTmp := map[interface{}]int{}
	for _, s := range list {
		if _, ok := mapTmp[s]; !ok {
			mapTmp[s] = 1
		}
	}
	var keys []interface{}
	for k := range mapTmp {
		keys = append(keys, k)
	}
	return keys
}

func MapTake(list []interface{}, limit int) []interface{} {
	res := []interface{}{}
	if limit > len(list) {
		limit = len(list)
	}
	for i := 0; i < limit; i++ {
		res = append(res, list[i])
	}
	return res
}

func ConverIntStringToList(eid string) (result string) {
	for i, e := range strings.Split(eid, ",") {
		v, err := strconv.Atoi(e)
		if err != nil {
			log.Debug(err.Error())
		} else {
			if i == 0 {
				result = fmt.Sprintf("%d", v)
			} else {
				result = fmt.Sprintf("%s,%d", result, v)
			}
		}
	}
	return
}
