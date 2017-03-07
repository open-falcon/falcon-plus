package pflag

import (
	"os"

	"github.com/spf13/pflag"
)

// Prints default by pflag and exit the application with 0
func PrintHelpAndExit0() {
	pflag.PrintDefaults()
	os.Exit(0)
}
