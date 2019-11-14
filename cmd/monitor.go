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

var Monitor = &cobra.Command{
	Use:   "monitor [Module ...]",
	Short: "Display an Open-Falcon module's log",
	Long: `
Display the log of the specified Open-Falcon module.
A module represents a single node in a cluster.
Modules:
  ` + strings.Join(g.AllModulesInOrder, " "),
	RunE: monitor,
}

func checkMonReq(name string) error {
	if !g.HasModule(name) {
		return fmt.Errorf("%s doesn't exist", name)
	}

	if !g.HasLogfile(name) {
		r := g.Rel(g.LogPath(name))
		return fmt.Errorf("expect logfile: %s", r)
	}

	return nil
}

func monitor(c *cobra.Command, args []string) error {
	if len(args) < 1 {
		return c.Usage()
	}
	var tailArgs []string = []string{"-f"}
	for _, moduleName := range args {
		if err := checkMonReq(moduleName); err != nil {
			return err
		}

		tailArgs = append(tailArgs, g.LogPath(moduleName))
	}
	cmd := exec.Command("tail", tailArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
