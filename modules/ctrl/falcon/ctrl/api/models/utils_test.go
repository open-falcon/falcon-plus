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
	"testing"
)

func TestMdiff(t *testing.T) {
	cases := []struct {
		src   []string
		dst   []string
		wanta []string
		wantd []string
	}{
		{src: []string{}, dst: []string{}, wanta: []string{}, wantd: []string{}},
		{src: []string{"1", "2"}, dst: []string{"3", "4"}, wanta: []string{"3", "4"}, wantd: []string{"1", "2"}},
		{src: []string{"1", "2"}, dst: []string{"2", "3"}, wanta: []string{"3"}, wantd: []string{"1"}},
	}
	for _, c := range cases {
		if gota, gotd := MdiffStr(c.src, c.dst); stringscmp(gota,
			c.wanta) != 0 || stringscmp(gotd, c.wantd) != 0 {
			t.Errorf("Mdiff(%v,%v) = %v, %v; want %v %v",
				c.src, c.dst, gota, gotd, c.wanta, c.wantd)
		}
	}
}

func Testintscmp64(t *testing.T) {
	cases := []struct {
		a    []int64
		b    []int64
		want bool
	}{
		{a: []int64{2, 3, 4}, b: []int64{3, 4, 2}, want: true},
		{a: []int64{2, 3, 4, 4}, b: []int64{3, 4, 4, 2}, want: true},
		{a: []int64{2, 3, 4, 4, 5}, b: []int64{3, 4, 4, 2, 6}, want: false},
		{a: []int64{2, 3, 4, 6}, b: []int64{3, 4, 2}, want: false},
	}
	for _, c := range cases {
		if got := intscmp64(c.a, c.b); (got == 0) != c.want {
			t.Errorf("intscmp64(%v,%v) = %v; want %v",
				c.a, c.b, got, c.want)
		}
	}
}
