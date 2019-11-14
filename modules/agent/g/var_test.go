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
	"path/filepath"
	"runtime"
	"testing"
)

func TestInitRootDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRootDir()
			if getCurrentPath() != Root {
				t.Errorf("Root: [%v], actually: [%v]", getCurrentPath(), Root)
			}
		})
	}
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filepath.Dir(filename))
}
