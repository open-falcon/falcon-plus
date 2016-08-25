package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/open-falcon/falcon-plus/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "open-falcon",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Commands is the mapping of all the available Consul commands.
//var Commands map[string]cli.CommandFactory
//var startCommand cli.Command

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.newAppName.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.AddCommand(cmd.Start)
	RootCmd.AddCommand(cmd.Stop)
	RootCmd.AddCommand(cmd.Restart)
	RootCmd.AddCommand(cmd.Check)
	RootCmd.AddCommand(cmd.Monitor)
	RootCmd.AddCommand(cmd.Reload)
	cmd.Start.Flags().BoolVar(&cmd.PreqOrderFlag, "preq-order", false, "start modules in the order of prerequisites")
	cmd.Start.Flags().BoolVar(&cmd.ConsoleOutputFlag, "console-output", false, "print the module's output to the console")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" { // enable ability to specify config file via flag
	//	viper.SetConfigFile(cfgFile)
	//}

	viper.SetConfigName(".newAppName") // name of config file (without extension)
	viper.AddConfigPath("$HOME")       // adding home directory as first search path
	viper.AutomaticEnv()               // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	log.SetOutput(ioutil.Discard)

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	//cli := &cli.CLI{
	//	Args:     args,
	//	Commands: Commands,
	//	//	HelpFunc: cli.BasicHelpFunc("consul"),
	//}
	//exitCode, err := cli.Run()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
	//	return 1
	//}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
