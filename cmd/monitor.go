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

func monitor(c *cobra.Command, args []string) error {
	if len(args) != 1 {
		return c.Usage()
	}
	moduleName := args[0]
	if !g.HasModule(moduleName) {
		return fmt.Errorf("%s doesn't exist\n", moduleName)
	}

	logPath := g.LogPath(moduleName)
	cmd := exec.Command("tail", "-f", logPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
