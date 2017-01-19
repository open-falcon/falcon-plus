package cmd

import (
	"fmt"
	"strings"

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
	if len(args) == 0 {
		return c.Usage()
	}
	if (len(args) == 1) && (args[0] == "all") {
		args = g.GetModuleArgsInOrder(g.AllModulesInOrder)
	} else {
		for _, moduleName := range args {
			err := g.ModuleExists(moduleName)
			if err != nil {
				fmt.Println(err)
				fmt.Println("** restart failed **")
				return nil //g.Command_EX_ERR
			}
		}
	}

	stop(c, args)
	start(c, args)
	return nil //g.Command_EX_OK
}
