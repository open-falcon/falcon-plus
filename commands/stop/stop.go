package stop

import (
	"fmt"
	"github.com/open-falcon/open-falcon/g"
	"github.com/mitchellh/cli"
	"os"
	"os/exec"
	"strings"
	"time"
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
				fmt.Println("** stop failed **")
				return g.Command_EX_ERR
			}
		}
	}
	for _, moduleName := range args {
		moduleStatus := g.CheckModuleStatus(moduleName)

		if moduleStatus == g.ModuleExistentNotRunning {
			// Skip stopping if the module is stopped
			continue
		}

		fmt.Print("Stopping [", g.ModuleApps[moduleName], "] ")

		pidStr, _ := g.CheckModulePid(g.ModuleApps[moduleName])

		cmd := exec.Command("kill", "-9", pidStr)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		fmt.Println("with PID [", pidStr, "]...successfully!!")
		time.Sleep(1 * time.Second)

		moduleStatus = g.CheckModuleStatus(moduleName)
		if moduleStatus == g.ModuleRunning {
			fmt.Println("** stop failed **")
			return g.Command_EX_ERR
		}
	}
	return g.Command_EX_OK
}

func (c *Command) Synopsis() string {
	return "Stop Open-Falcon modules"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon stop [Module ...]

  Stop the specified Open-Falcon modules.
  A module represents a single node in a cluster.

Modules:

  ` + "all " + strings.Join(g.AllModulesInOrder, " ")
	return strings.TrimSpace(helpText)
}
