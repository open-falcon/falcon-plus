package main

import (
	"fmt"
	"os"

	"github.com/open-falcon/falcon-plus/cmd"
	"github.com/spf13/cobra"
)

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
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
