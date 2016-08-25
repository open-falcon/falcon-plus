package cmd

import (
	"fmt"
	"strings"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

var Check = &cobra.Command{
	Use:   "status [Module ...]",
	Short: "Check the status of Open-Falcon modules",
	Long: `
Check if the specified Open-Falcon modules are running, not running or noexistent.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: check,
}

func check(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return c.Usage()
	}
	if (len(args) == 1) && (args[0] == "all") {
		args = g.AllModulesInOrder
	} else {
		for _, moduleName := range args {
			err := g.ModuleExists(moduleName)
			if err != nil {
				fmt.Println(err)
				fmt.Println("** status failed **")
				return nil //g.Command_EX_ERR
			}
		}
	}
	for _, moduleName := range args {
		g.CheckModuleStatus(moduleName)
	}
	return nil //g.Command_EX_OK
}
