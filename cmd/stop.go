package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/spf13/cobra"
)

var Stop = &cobra.Command{
	Use:   "stop [Module ...]",
	Short: "Stop Open-Falcon modules",
	Long: `
Stop the specified Open-Falcon modules.
A module represents a single node in a cluster.
Modules:
  ` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: stop,
}

func stop(c *cobra.Command, args []string) error {
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
				fmt.Println("** stop failed **")
				return nil //g.Command_EX_ERR
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
			return nil //g.Command_EX_ERR
		}
	}
	return nil //g.Command_EX_OK
}
