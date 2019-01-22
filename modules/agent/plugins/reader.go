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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
)

// key: sys/ntp/60_ntp.py
func ListPlugins(relativePath string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if relativePath == "" {
		return ret
	}

	//解析参数
	var args string
	re := regexp.MustCompile(`(.*)\((.*)\)`)
	relPathWithArgs := re.FindAllStringSubmatch(relativePath, -1)
	if relPathWithArgs == nil {
		relativePath = relativePath
		args = ""
	} else {
		relativePath = relPathWithArgs[0][1]
		args = relPathWithArgs[0][2]
	}

	path := filepath.Join(g.Config().Plugin.Dir, relativePath)

	//处理路径为脚本的情况
	if file.IsFile(path) {
		dir, fileName := filepath.Split(path)
		arr := strings.Split(fileName, "_")
		var cycle int
		var err error
		cycle, err = strconv.Atoi(arr[0])
		if err == nil {
			fi, _ := os.Stat(path)
			plugin := &Plugin{FilePath: relativePath, MTime: fi.ModTime().Unix(), Cycle: cycle, Args: args}
			ret[dir+"("+args+")"] = plugin
			return ret
		}
	}

	if !file.IsExist(path) || file.IsFile(path) {
		return ret
	}

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("can not list files under", path)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		args = ""
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
		plugin := &Plugin{FilePath: fpath, MTime: f.ModTime().Unix(), Cycle: cycle, Args: args}
		ret[fpath+"("+args+")"] = plugin
	}
	return ret
}
