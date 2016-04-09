package start

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
		args = g.GetModuleArgsInOrder(g.AllModulesInOrder)
	} else {
		for _, moduleName := range args {
			err := g.ModuleExists(moduleName)
			if err != nil {
				fmt.Println(err)
				fmt.Println("** start failed **")
				return g.Command_EX_ERR
			}
		}
		args = g.GetModuleArgsInOrder(args)
	}
	for _, moduleName := range args {
		moduleStatus := g.CheckModuleStatus(moduleName)

		if moduleStatus == g.ModuleExistentNotRunning {
			fmt.Print("Starting [", g.ModuleApps[moduleName], "]...")
			cmdArgs, err := g.GetConfFileArgs(g.ModuleConfs[moduleName])
			if err != nil {
				fmt.Println(err)
				return g.Command_EX_ERR
			}
			// fe workaround
			if moduleName == "fe" {
				os.Chdir("bin/fe")
				cmd := exec.Command("./control", "start")
				dir, _ := os.Getwd()
				cmd.Dir = dir
				cmd.Start()
				fmt.Println("successfully!!")
				time.Sleep(1 * time.Second)
				moduleStatus = g.CheckModuleStatus(moduleName)
				if moduleStatus == g.ModuleExistentNotRunning {
					return g.Command_EX_ERR
				}
				os.Chdir("../../")
				continue
			}
			logPath := g.LogDir + "/" + moduleName + ".log"
			LogOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("Error in opening file:", err)
				return g.Command_EX_ERR
			}
			defer LogOutput.Close()

			cmd := exec.Command(g.ModuleBins[moduleName], cmdArgs...)
			cmd.Stdout = LogOutput
			cmd.Stderr = LogOutput
			dir, _ := os.Getwd()
			cmd.Dir = dir
			cmd.Start()
			fmt.Println("successfully!!")
			time.Sleep(1 * time.Second)
			moduleStatus = g.CheckModuleStatus(moduleName)
			if moduleStatus == g.ModuleExistentNotRunning {
				fmt.Println("** start failed **")
				return g.Command_EX_ERR
			}
		}
		// Skip starting if the module is already running
	}
	return g.Command_EX_OK
}

func (c *Command) Synopsis() string {
	return "Start Open-Falcon modules"
}

func (c *Command) Help() string {
	helpText := `
Usage: open-falcon start [Module ...]

  Start the specified Open-Falcon modules and run until a stop command is received.
  A module represents a single node in a cluster.

Modules:

  ` + "all " + strings.Join(g.AllModulesInOrder, " ")
	return strings.TrimSpace(helpText)
}
