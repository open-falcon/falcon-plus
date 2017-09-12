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

package cron

import (
	"testing"
)

func Test_expressionValid(t *testing.T) {

	expressionMap := map[string]bool{
		// true
		"1210":                                     true,
		"$#":                                       true,
		"$(cpu.busy)":                              true,
		"$(cpu.busy)+$(cpu.idle)-$(cpu.nice)":      true,
		"$(cpu.busy)>=80":                          true,
		"($(cpu.busy)+$(cpu.idle)-$(cpu.nice))>80": true,
		"$(qps/module=judge,project=falcon)":       true,
		"($(cpu.idle)+$(cpu.busy))=100":            true,

		// false
		"$((cpu.busy)":                     false,
		"$(cpu.idle)+$(cpu.busy)>40":       false,
		"($(cpu.idle)+$(cpu.busy)-60)>100": false,
	}

	for key, val := range expressionMap {
		if st := expressionValid(key); st != val {
			t.Errorf("func expressionValid() failure")
		}
	}
}
