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
	"testing"
)

var testCases4map2string = []struct {
	tags   map[string]string
	expect string
}{
	{map[string]string{"1": "1"}, "1=1"},
	{map[string]string{"1": "1", "2": "2"}, "1=1,2=2"},
	{map[string]string{"1": "1", "2": "2", "0": "0"}, "0=0,1=1,2=2"},
}

var testCases4string2map = []struct {
	tags   string
	expect map[string]string
}{
	{"1=1", map[string]string{"1": "1"}},
	{"1=1,2=2", map[string]string{"1": "1", "2": "2"}},
	{"0=0,1=1,2=2", map[string]string{"1": "1", "2": "2", "0": "0"}},
	{"0,1=1,2=2", map[string]string{"1": "1", "2": "2"}},
	{"0=,1=1,2=2", map[string]string{"0": "", "1": "1", "2": "2"}},
	{"=0,1=1,2=2", map[string]string{"": "0", "1": "1", "2": "2"}},
	{"0=0, 1=1, 2=2", map[string]string{"0": "0", "1": "1", "2": "2"}},
}

func origSortedTags(tags map[string]string) string {
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

func origDictedTagstring(s string) map[string]string {
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

func Test_SortedTags(t *testing.T) {
	for _, testCase := range testCases4map2string {
		if r := SortedTags(testCase.tags); r != testCase.expect || SortedTags(testCase.tags) != origSortedTags(testCase.tags) {
			t.Errorf("expect %v, got %v\n", testCase.expect, r)
		}
	}
}
func Test_DictedTagstring(t *testing.T) {
	for _, testCase := range testCases4string2map {
		r := DictedTagstring(testCase.tags)
		origR := origDictedTagstring(testCase.tags)

		if len(r) != len(testCase.expect) || len(r) != len(origR) {
			t.FailNow()
		}
		for k, v := range r {
			if expectV, exist := testCase.expect[k]; !exist || v != expectV || v != origR[k] {
				t.Errorf("expect %v, got %v\n", testCase.expect, r)
			}
		}
	}
}

func Benchmark_SortedTags_1pair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1"})
	}
}

func Benchmark_SortedTags_1pair_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1"})
	}
}

func Benchmark_SortedTags_2pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1", "2": "2"})
	}
}

func Benchmark_SortedTags_2pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1", "2": "2"})
	}
}

func Benchmark_SortedTags_3pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1", "2": "2", "3": "3"})
	}
}

func Benchmark_SortedTags_3pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1", "2": "2", "3": "3"})
	}
}

func Benchmark_SortedTags_4pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4"})
	}
}

func Benchmark_SortedTags_4pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4"})
	}
}

func Benchmark_SortedTags_5pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5"})
	}
}

func Benchmark_SortedTags_5pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5"})
	}
}

func Benchmark_SortedTags_6pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "0": "0"})
	}
}

func Benchmark_SortedTags_6pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origSortedTags(map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "0": "0"})
	}
}

func Benchmark_DictedTagstring_1pair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1")
	}
}

func Benchmark_DictedTagstring_1pair_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1")
	}
}

func Benchmark_DictedTagstring_2pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1,2=2")
	}
}

func Benchmark_DictedTagstring_2pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1,2=2")
	}
}

func Benchmark_DictedTagstring_3pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1,2=2,3=3")
	}
}

func Benchmark_DictedTagstring_3pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1,2=2,3=3")
	}
}

func Benchmark_DictedTagstring_4pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1,2=2,3=3,4=4")
	}
}

func Benchmark_DictedTagstring_4pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1,2=2,3=3,4=4")
	}
}

func Benchmark_DictedTagstring_5pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1,2=2,3=3,4=4,5=5")
	}
}

func Benchmark_DictedTagstring_5pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1,2=2,3=3,4=4,5=5")
	}
}

func Benchmark_DictedTagstring_6pairs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DictedTagstring("1=1,2=2,3=3,4=4,5=5,6=6")
	}
}

func Benchmark_DictedTagstring_6pairs_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origDictedTagstring("1=1,2=2,3=3,4=4,5=5,6=6")
	}
}
