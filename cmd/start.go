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
	RunE:          start,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func cmdArgs(name string) []string {
	return []string{"-c", g.Cfg(name)}
}

func openLogFile(name string) (*os.File, error) {
	logDir := g.LogDir(name)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logPath := g.LogPath(name)
	logOutput, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logOutput, nil
}

func checkReq(name string) error {
	if err := g.HasModule(name); err != nil {
		return err
	}

	if err := g.HasCfg(name); err != nil {
		return err
	}

	return nil
}

func isStarted(name string) bool {
	ticker := time.NewTicker(time.Millisecond * 100)

	for {
		select {
		case <-ticker.C:
			if g.IsRunning(name) {
				return true
			}
		case <-time.After(time.Second):
			return false
		}
	}
}

func start(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return c.Usage()
	}
	//if (len(args) == 1) && (args[0] == "all") {
	//	args = g.GetModuleArgsInOrder(g.Order)
	//} else {
	//	for _, moduleName := range args {
	//		err := g.HasModule(moduleName)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//	args = g.GetModuleArgsInOrder(args)
	//}
	g.PreqOrder(args)

	for _, moduleName := range args {
		// Skip starting if the module is already running
		if g.IsRunning(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		if err := checkReq(moduleName); err != nil {
			return err
		}

		logOutput, err := openLogFile(moduleName)
		if err != nil {
			return err
		}
		defer logOutput.Close()

		cmd := exec.Command(g.Bin(moduleName), cmdArgs(moduleName)...)
		cmd.Stdout = logOutput
		cmd.Stderr = logOutput
		dir, _ := os.Getwd()
		cmd.Dir = dir
		cmd.Start()

		if isStarted(moduleName) {
			fmt.Print("[", g.ModuleApps[moduleName], "] ", g.Pid(moduleName), "\n")
			continue
		}

		return fmt.Errorf("[%s] Failed to start", g.ModuleApps[moduleName])
	}
	return nil
}
