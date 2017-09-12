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
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

var Restart = &cobra.Command{
	Use:   "restart [Module ...]",
	Short: "Restart Open-Falcon modules",
	Long: `
Restart the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: restart,
}

func restart(c *cobra.Command, args []string) error {
	args = g.RmDup(args)

	if len(args) == 0 {
		args = g.AllModulesInOrder
	}

	for _, moduleName := range args {
		if err := stop(c, []string{moduleName}); err != nil {
			return err
		}
		if strings.Contains(moduleName, "graph") {
			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(1 * time.Second)
		}
		if err := start(c, []string{moduleName}); err != nil {
			return err
		}
	}
	return nil
}
