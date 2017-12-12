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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"io/ioutil"
)

func HasLogfile(name string) bool {
	if _, err := os.Stat(LogPath(name)); err != nil {
		return false
	}
	return true
}

func PreqOrder(moduleArgs []string) []string {
	if len(moduleArgs) == 0 {
		return []string{}
	}

	var modulesInOrder []string

	// get arguments which are found in the order
	for _, nameOrder := range AllModulesInOrder {
		for _, nameArg := range moduleArgs {
			if nameOrder == nameArg {
				modulesInOrder = append(modulesInOrder, nameOrder)
			}
		}
	}
	// get arguments which are not found in the order
	for _, nameArg := range moduleArgs {
		end := 0
		for _, nameOrder := range modulesInOrder {
			if nameOrder == nameArg {
				break
			}
			end++
		}
		if end == len(modulesInOrder) {
			modulesInOrder = append(modulesInOrder, nameArg)
		}
	}
	return modulesInOrder
}

func Rel(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// filepath.Abs() returns an error only when os.Getwd() returns an error;
	abs, _ := filepath.Abs(p)

	r, err := filepath.Rel(wd, abs)
	if err != nil {
		return ""
	}

	return r
}

func HasCfg(name string) bool {
	if _, err := os.Stat(Cfg(name)); err != nil {
		return false
	}
	return true
}

func HasModule(name string) bool {
	return Modules[name]
}

func setPid(name string) {
	output, _ := exec.Command("pgrep", "-f", ModuleApps[name]).Output()
	pidStr := strings.TrimSpace(string(output))
	//Write the pid in the file
	pid := pidStr + "\n"
	pid_file := ModuleApps[name][7 : ]+"/logs/"+ModuleApps[name][7 : ]+".pid"
	ioutil.WriteFile(pid_file, []byte(pid), 0644)

	PidOf[name] = pidStr
}

func Pid(name string) string {
	if PidOf[name] == "<NOT SET>" {
		setPid(name)
	}
	return PidOf[name]
}

func IsRunning(name string) bool {
	setPid(name)
	return Pid(name) != ""
}

func RmDup(args []string) []string {
	if len(args) == 0 {
		return []string{}
	}
	if len(args) == 1 {
		return args
	}

	ret := []string{}
	isDup := make(map[string]bool)
	for _, arg := range args {
		if isDup[arg] == true {
			continue
		}
		ret = append(ret, arg)
		isDup[arg] = true
	}
	return ret
}
