package g

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/toolkits/file"
)

func configExists(cfg string) bool {
	if !file.IsExist(cfg) {
		return false
	}
	return true
}

var regexpReplaceCurrentFolder, _ = regexp.Compile("^\\.")

func ConfFileArgs(cfg string) ([]string, error) {
	if !file.IsExist(cfg) {
		return nil, fmt.Errorf("expect config file: %s\n", cfg)
	}
	return []string{"-c", cfg}, nil
}

func ModuleExists(name string) error {
	if Modules[name] {
		return nil
	}
	return fmt.Errorf("This module doesn't exist: %s", name)
}

func CheckModulePid(name string) (string, error) {
	output, err := exec.Command("pgrep", name).Output()
	if err != nil {
		return "", err
	}
	pidStr := strings.TrimSpace(string(output))
	return pidStr, nil
}

func PreqOrder(moduleArgs []string) []string {
	var modulesInOrder []string

	// get arguments which are found in the order
	for _, nameOrder := range AllModulesInOrder {
		for _, nameArg := range moduleArgs {
			if nameOrder == nameArg {
				modulesInOrder = append(modulesInOrder, nameOrder)
			}
		}
	}
	// get arguments which are not found in the order
	for _, nameArg := range moduleArgs {
		end := 0
		for _, nameOrder := range modulesInOrder {
			if nameOrder == nameArg {
				break
			}
			end++
		}
		if end == len(modulesInOrder) {
			modulesInOrder = append(modulesInOrder, nameArg)
		}
	}
	return modulesInOrder
}

func CheckModuleStatus(name string) int {
	fmt.Print("Checking status [", ModuleApps[name], "]...")

	pidStr, err := CheckModulePid(ModuleApps[name])
	if err != nil {
		fmt.Println("not running!!")
		return ModuleExistentNotRunning
	}

	fmt.Println("running with PID [", pidStr, "]!!")
	return ModuleRunning
}

func Rel(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// filepath.Abs() returns an error only when os.Getwd() returns an error;
	abs, _ := filepath.Abs(p)

	r, err := filepath.Rel(wd, abs)
	if err != nil {
		return ""
	}

	return r
}

func HasCfg(name string) bool {
	if _, err := os.Stat(Cfg(name)); err != nil {
		return false
	}
	return true
}

func HasModule(name string) bool {
	if Modules[name] {
		return true
	}
	return false
}

func setPid(name string) {
	output, _ := exec.Command("pgrep", "-f", ModuleApps[name]).Output()
	pidStr := strings.TrimSpace(string(output))
	PidOf[name] = pidStr
}

func Pid(name string) string {
	if PidOf[name] == "<NOT SET>" {
		setPid(name)
	}
	return PidOf[name]
}

func IsRunning(name string) bool {
	setPid(name)
	if Pid(name) == "" {
		return false
	}
	return true
}
