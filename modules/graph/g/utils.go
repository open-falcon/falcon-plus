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

package g

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/toolkits/file"
)

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
func RrdFileName(baseDir string, md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s/%s/%s_%s_%d.rrd", baseDir, md5[0:2], md5, dsType, step)
}

// rrd文件是否存在
func IsRrdFileExist(filename string) bool {
	return file.IsExist(filename)
}

// 生成rrd缓存数据的key
func FormRrdCacheKey(md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s_%s_%d", md5, dsType, step)
}
func SplitRrdCacheKey(ckey string) (md5 string, dsType string, step int, err error) {
	ckey_slice := strings.Split(ckey, "_")
	if len(ckey_slice) != 3 {
		err = fmt.Errorf("bad rrd cache key: %s", ckey)
		return
	}

	md5 = ckey_slice[0]
	dsType = ckey_slice[1]
	stepInt64, err := strconv.ParseInt(ckey_slice[2], 10, 32)
	if err != nil {
		return
	}
	step = int(stepInt64)

	// return
	err = nil
	return
}

// 判断是否为有效字符串(不包含指定字符和多字节字符)
func IsValidString(str string) bool {

	r := []rune(str)
	// 多字节字符
	if len(r) != len(str) {
		return false
	}

	for _, t := range r {
		switch t {
		case '\r':
			return false
		case '\n':
			return false
		case '\'':
			return false
		case '"':
			return false
		case '>':
			return false
		case '\032':
			return false
		default:
			// 不可打印字符
			if !unicode.IsPrint(t) {
				return false
			}
		}
	}
	return true
}
