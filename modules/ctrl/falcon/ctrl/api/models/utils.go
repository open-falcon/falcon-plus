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

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Log struct {
	Id       int64
	Module   int64
	ModuleId int64
	UserId   int64
	Action   int64
	Data     string
	Time     time.Time
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

type cache struct {
	enable bool
	data   map[int64]interface{}
}

func (c *cache) set(id int64, p interface{}) {
	if c.enable {
		c.data[id] = p
	}
}

func (c *cache) get(id int64) interface{} {
	return c.data[id]
}

func (c *cache) del(id int64) {
	if c.enable {
		delete(c.data, id)
	}
}

func DbLog(o orm.Ormer, uid, module, module_id, action int64, data string) {
	log := &Log{
		UserId:   uid,
		Module:   module,
		ModuleId: module_id,
		Action:   action,
		Data:     data,
	}
	o.Insert(log)
}

func array2sql(array []int64) string {
	var ret string
	if len(array) == 0 {
		return "()"
	}

	for i := 0; i < len(array); i++ {
		ret += fmt.Sprintf("%d,", array[i])
	}
	return fmt.Sprintf("(%s)", ret[:len(ret)-1])
}

func stringscmp(a, b []string) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := 0; i < len(a); i++ {
		if ret = strings.Compare(a[i], b[i]); ret != 0 {
			return
		}
	}
	return
}

func intscmp64(a, b []int64) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}

	_a := make([]int, len(a))
	for i := 0; i < len(_a); i++ {
		_a[i] = int(a[i])
	}

	_b := make([]int, len(b))
	for i := 0; i < len(_b); i++ {
		_b[i] = int(b[i])
	}

	sort.Ints(_a)
	sort.Ints(_b)

	for i := 0; i < len(_a); i++ {
		if ret = _a[i] - _b[i]; ret != 0 {
			return
		}
	}
	return
}

func intscmp(a, b []int) (ret int) {
	if ret = len(a) - len(b); ret != 0 {
		return
	}
	sort.Ints(a)
	sort.Ints(b)
	for i := 0; i < len(a); i++ {
		if ret = a[i] - b[i]; ret != 0 {
			return
		}
	}
	return
}

func jsonStr(i interface{}) string {
	if ret, err := json.Marshal(i); err != nil {
		return ""
	} else {
		return string(ret)
	}
}

func MdiffStr(src, dst []string) (add, del []string) {
	_src := make(map[string]bool)
	_dst := make(map[string]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}
func MdiffInt(src, dst []int64) (add, del []int64) {
	_src := make(map[int64]bool)
	_dst := make(map[int64]bool)
	for _, v := range src {
		_src[v] = true
	}
	for _, v := range dst {
		_dst[v] = true
	}
	for k, _ := range _src {
		if !_dst[k] {
			del = append(del, k)
		}
	}
	for k, _ := range _dst {
		if !_src[k] {
			add = append(add, k)
		}
	}
	return
}

func GetIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}
