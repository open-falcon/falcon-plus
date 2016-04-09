package g

import (
	"fmt"
	"github.com/toolkits/file"
	"os/exec"
	"strings"
)

func configExists(cfg string) bool {
	if !file.IsExist(cfg) {
		return false
	}
	return true
}

func GetConfFileArgs(cfg string) ([]string, error) {
	if !configExists(cfg) {
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

func GetModuleArgsInOrder(moduleArgs []string) []string {
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
