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
	"strconv"
	"strings"
)

func ReadableFloat(raw float64) string {
	val := strconv.FormatFloat(raw, 'f', 5, 64)
	if strings.Contains(val, ".") {
		val = strings.TrimRight(val, "0")
		val = strings.TrimRight(val, ".")
	}

	return val
}
