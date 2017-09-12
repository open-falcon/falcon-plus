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

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

var Stop = &cobra.Command{
	Use:   "stop [Module ...]",
	Short: "Stop Open-Falcon modules",
	Long: `
Stop the specified Open-Falcon modules.
A module represents a single node in a cluster.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: stop,
}

func stop(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist", moduleName)
		}

		if !g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
			continue
		}

		cmd := exec.Command("kill", "-TERM", g.Pid(moduleName))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err == nil {
			fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
			continue
		}
		return err
	}
	return nil
}
