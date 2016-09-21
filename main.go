package main

import (
	"fmt"
	"os"

	"github.com/open-falcon/falcon-plus/cmd"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var versionFlag bool

var RootCmd = &cobra.Command{
	Use: "open-falcon",
}

func init() {
	RootCmd.AddCommand(cmd.Start)
	RootCmd.AddCommand(cmd.Stop)
	RootCmd.AddCommand(cmd.Restart)
	RootCmd.AddCommand(cmd.Check)
	RootCmd.AddCommand(cmd.Monitor)
	RootCmd.AddCommand(cmd.Reload)
	cmd.Start.Flags().BoolVar(&cmd.PreqOrderFlag, "preq-order", false, "start modules in the order of prerequisites")
	cmd.Start.Flags().BoolVar(&cmd.ConsoleOutputFlag, "console-output", false, "print the module's output to the console")
	flag.BoolVarP(&versionFlag, "version", "v", false, "show version")
	flag.Parse()
}

func main() {
	if versionFlag {
		fmt.Printf("Open-Falcon version %s, build %s\n", Version, GitCommit)
		os.Exit(0)
	}
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
