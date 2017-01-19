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

var Start = &cobra.Command{
	Use:   "start [Module ...]",
	Short: "Start Open-Falcon modules",
	Long: `
Start the specified Open-Falcon modules and run until a stop command is received.
A module represents a single node in a cluster.
Modules:
	` + "all " + strings.Join(g.AllModulesInOrder, " "),
	RunE: start,
}

func start(c *cobra.Command, args []string) error {
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
				fmt.Println("** start failed **")
				return nil //g.Command_EX_ERR
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
				return nil //g.Command_EX_ERR
			}

			logPath := "./" + moduleName + "/" + g.LogDir + "/" + moduleName + ".log"
			LogOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("Error in opening file:", err)
				return nil //g.Command_EX_ERR
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
				return nil //g.Command_EX_ERR
			}
		}
		// Skip starting if the module is already running
	}
	return nil //g.Command_EX_OK
}
