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
	"sort"
	"strings"
)

func SortedTags(tags map[string]string) string {
	if tags == nil {
		return ""
	}

	size := len(tags)

	if size == 0 {
		return ""
	}

	if size == 1 {
		for k, v := range tags {
			return fmt.Sprintf("%s=%s", k, v)
		}
	}

	keys := make([]string, size)
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	ret := make([]string, size)
	for j, key := range keys {
		ret[j] = fmt.Sprintf("%s=%s", key, tags[key])
	}

	return strings.Join(ret, ",")
}

func DictedTagstring(s string) map[string]string {
	if s == "" {
		return map[string]string{}
	}
	s = strings.Replace(s, " ", "", -1)

	tag_dict := make(map[string]string)
	tags := strings.Split(s, ",")
	for _, tag := range tags {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {
			tag_dict[tag_pair[0]] = tag_pair[1]
		}
	}
	return tag_dict
}

func SplitTagsString(s string) (err error, tags map[string]string) {
	err = nil
	tags = make(map[string]string)

	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		return
	}

	tagSlice := strings.Split(s, ",")
	for _, tag := range tagSlice {
		tag_pair := strings.SplitN(tag, "=", 2)
		if len(tag_pair) == 2 {
			tags[tag_pair[0]] = tag_pair[1]
		} else {
			err = fmt.Errorf("bad tag %s", tag)
			return
		}
	}

	return
}
