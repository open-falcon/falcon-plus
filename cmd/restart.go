package cmd

import (
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
	args = g.RmDup(args)

	for _, moduleName := range args {
		if err := stop(c, []string{moduleName}); err != nil {
			return err
		}
		if err := start(c, []string{moduleName}); err != nil {
			return err
		}
	}
	return nil
}
