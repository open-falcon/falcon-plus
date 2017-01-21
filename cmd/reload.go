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
