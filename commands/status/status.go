package status

import (
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/open-falcon/open-falcon/g"
	"strings"
)

// Command is a Command implementation that runs a Consul agent.
// The command will not end unless a shutdown message is sent on the
// ShutdownCh. If two messages are sent on the ShutdownCh it will forcibly
// exit.
type Command struct {
	Revision          string
	Version           string
	VersionPrerelease string
	Ui                cli.Ui
}

func (c *Command) Run(args []string) int {
	if len(args) == 0 {
		return cli.RunResultHelp
	}
	if (len(args) == 1) && (args[0] == "all") {
		args = g.AllModulesInOrder
	} else {
		for _, moduleName := range args {
			err := g.ModuleExists(moduleName)
			if err != nil {
				fmt.Println(err)
				fmt.Println("** status failed **")
				return g.Command_EX_ERR
			}
		}
	}
	for _, moduleName := range args {
		g.CheckModuleStatus(moduleName)
	}
	return g.Command_EX_OK
}

func (c *Command) Synopsis() string {
	return "Check the status of Open-Falcon modules"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon status [Modules ...]

  Check if the specified Open-Falcon modules are running, not running or noexistent.

Modules:

  ` + "all " + strings.Join(g.AllModulesInOrder, " ")
	return strings.TrimSpace(helpText)
}
