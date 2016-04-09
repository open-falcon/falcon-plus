package restart

import (
	"fmt"
	"github.com/open-falcon/open-falcon/commands/start"
	"github.com/open-falcon/open-falcon/commands/stop"
	"github.com/open-falcon/open-falcon/g"
	"github.com/mitchellh/cli"
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
		args = g.GetModuleArgsInOrder(g.AllModulesInOrder)
	} else {
		for _, moduleName := range args {
			err := g.ModuleExists(moduleName)
			if err != nil {
				fmt.Println(err)
				fmt.Println("** restart failed **")
				return g.Command_EX_ERR
			}
		}
	}

	var stopCmd stop.Command
	var startCmd start.Command
	stopCmd.Run(args)
	startCmd.Run(args)
	return g.Command_EX_OK
}

func (c *Command) Synopsis() string {
	return "Restart Open-Falcon modules"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon restart [Module ...]

  Restart the specified Open-Falcon modules and run until a stop command is received.
  A module represents a single node in a cluster.

Modules:

  ` + "all " + strings.Join(g.AllModulesInOrder, " ")
	return strings.TrimSpace(helpText)
}
