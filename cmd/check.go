package cmd

import (
	"fmt"
	"strings"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

var Check = &cobra.Command{
	Use:   "check [Module ...]",
	Short: "Check the status of Open-Falcon modules",
	Long: `
Check if the specified Open-Falcon modules are running.

Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: check,
}

func check(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return c.Usage()
	}
	for _, moduleName := range args {
		if !g.HasModule(moduleName) {
			return fmt.Errorf("%s doesn't exist\n", moduleName)
		}

		if g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
		} else {
			fmt.Print("[", g.ModuleApps[moduleName], "] down\n")
		}
	}
	return nil
}
