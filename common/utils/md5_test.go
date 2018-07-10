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
	"crypto/md5"
	"fmt"
	"io"
	"testing"
)

func origMd5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Test_Md5(t *testing.T) {
	if Md5("1234567890123") != origMd5("1234567890123") {
		t.Error("not expect")
	}
}

func Benchmark_Md5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Md5("1234567890123")
	}
}

func Benchmark_Md5_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origMd5("1234567890123")
	}
}
