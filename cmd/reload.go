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

import "github.com/spf13/cobra"

var Reload = &cobra.Command{
	Use:   "reload [Module ...]",
	Short: "Reload an Open-Falcon module's configuration file",
	Long: `
Reload the configuration file of the specified Open-Falcon module.
A module represents a single node in a cluster.
Modules:
  `,
	RunE: reload,
}

func reload(c *cobra.Command, args []string) error {
	if len(args) != 1 {
		return c.Usage()
	}
	return nil
}
