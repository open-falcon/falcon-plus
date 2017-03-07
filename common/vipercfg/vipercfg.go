package vipercfg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/chyeh/viper"
	owlPflag "github.com/open-falcon/falcon-plus/common/pflag"
	"github.com/spf13/pflag"
)

// Gets called before pflag.Parse()
type PflagDefiner func()

// Default Pflag definer of OWL
func OwlDefaultPflagDefiner() {
	pflag.StringP("config", "c", "cfg.json", "configuration file")
	pflag.BoolP("version", "v", false, "show version")
	pflag.Bool("check", false, "check collector")
	pflag.BoolP("help", "h", false, "usage")
	pflag.Bool("vg", false, "show version and git commit log")
}

// Used for parsing arguments of command line and loading file of configuration
//
// The method could be called by multiple times without reloading.
//
// LoadConfigFile() would call ParseCmd() to get the file path of configuration.
type ConfigLoader struct {
	// The definer of pflag
	FlagDefiner PflagDefiner
	// If the value of string is true, the mapped function would be called
	TrueValueCallbacks map[string]func()

	cmdViper    *viper.Viper
	configViper *viper.Viper

	parseCmdError       error
	loadConfigFileError error
}

func NewOwlConfigLoader() *ConfigLoader {
	return &ConfigLoader{
		FlagDefiner: OwlDefaultPflagDefiner,
		TrueValueCallbacks: map[string]func(){
			"help": owlPflag.PrintHelpAndExit0,
		},
	}
}

// Loads configuration from command line, panic if there is error while parsing
func (loader *ConfigLoader) MustParseCmd() *viper.Viper {
	viper, err := loader.ParseCmd()
	if err != nil {
		panic(loader.parseCmdError)
	}

	return viper
}

// Loads configuration from command line
func (loader *ConfigLoader) ParseCmd() (*viper.Viper, error) {
	if loader.cmdViper != nil || loader.parseCmdError != nil {
		return loader.cmdViper, loader.parseCmdError
	}

	/**
	 * Loads arguments from command line
	 */
	loader.FlagDefiner()
	pflag.Parse()

	loader.cmdViper = viper.New()
	pflag.VisitAll(func(flag *pflag.Flag) {
		/**
		 * Skip loading for the rest of flags
		 */
		if loader.parseCmdError != nil {
			return
		}
		// :~)

		if err := loader.cmdViper.BindPFlag(flag.Name, pflag.Lookup(flag.Name)); err != nil {
			loader.cmdViper = nil
			loader.parseCmdError = err
		}
	})

	return loader.cmdViper, loader.parseCmdError
}

// Loads configuration from command line, panic if there is error while parsing JSON
func (loader *ConfigLoader) MustLoadConfigFile() *viper.Viper {
	viper, err := loader.LoadConfigFile()

	if err != nil {
		panic(err)
	}

	return viper
}

// Executes callback for true values
func (loader *ConfigLoader) ProcessTrueValueCallbacks() {
	cmdViper := loader.MustParseCmd()

	for name, callback := range loader.TrueValueCallbacks {
		if cmdViper.GetBool(name) {
			callback()
		}
	}
}

// Loads configuration from command line, panic if there is error
func (loader *ConfigLoader) LoadConfigFile() (*viper.Viper, error) {
	if loader.configViper != nil || loader.loadConfigFileError != nil {
		return loader.configViper, loader.loadConfigFileError
	}

	/**
	 * Loads file path of configuration from command line
	 */
	cmdViper, cmdError := loader.ParseCmd()
	if cmdError != nil {
		loader.loadConfigFileError = cmdError
		return loader.configViper, loader.loadConfigFileError
	}
	// :~)

	/**
	 * Loads properties from config file
	 */
	cfgPath := cmdViper.GetString("config")

	jsonConfigViper := viper.New()
	jsonConfigViper.SetConfigFile(cfgPath)

	err := jsonConfigViper.ReadInConfig()
	if err != nil {
		loader.loadConfigFileError = err
	} else {
		loader.configViper = jsonConfigViper
	}
	// :~)

	return loader.configViper, loader.loadConfigFileError
}

// Loads the configuration which combines both arguments of command line and
// configuration file.
//
// Priorities:
//
// 1. Arguments of command line, which overrides:
// 2. Properties of config file
func (loader *ConfigLoader) BuildFacadeConfig() (*viper.Viper, error) {
	cmdViper := loader.MustParseCmd()
	configFileViper, fileError := loader.LoadConfigFile()

	if fileError != nil {
		return nil, fileError
	}

	newViper := viper.New()

	/**
	 * Use command line overrides properties of configuration file
	 */
	for k, v := range configFileViper.AllSettings() {
		newViper.Set(k, v)
	}
	for k, v := range cmdViper.AllSettings() {
		newViper.Set(k, v)
	}
	// :~)

	return newViper, nil
}

// ====================
// Following functions is deprecated
// ====================

var viperLoader = NewOwlConfigLoader()

// Deprecated
func Load() {
	viperLoader.MustLoadConfigFile()
}

// Deprecated
func Parse() {
	viperLoader.MustParseCmd()
}

// Deprecated
func Bind() {
	viperLoader.MustParseCmd()
}

// Deprecated
func Config() *viper.Viper {
	viper, err := viperLoader.BuildFacadeConfig()

	if err != nil {
		log.Warnf("Cannot load configuration. Error: %v", err)
		return nil
	}

	return viper
}
