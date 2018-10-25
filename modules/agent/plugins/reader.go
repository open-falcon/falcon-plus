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

package plugins

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"path/filepath"
	"strconv"
	"strings"
)

// key: sys/ntp/60_ntp.py
func ListPlugins(relativePath string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if relativePath == "" {
		return ret
	}

	dir := filepath.Join(g.Config().Plugin.Dir, relativePath)

	if !file.IsExist(dir) || file.IsFile(dir) {
		return ret
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("can not list files under", dir)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $cycle_$xx
		var cycle int
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		fpath := filepath.Join(relativePath, filename)
		plugin := &Plugin{FilePath: fpath, MTime: f.ModTime().Unix(), Cycle: cycle}
		ret[fpath] = plugin
	}

	return ret
}
