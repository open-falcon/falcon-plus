package cron

import (
	"testing"
)

func Test_expressionValid(t *testing.T) {

	expressionMap := map[string]bool{
		// true
		"1210":                                     true,
		"$#":                                       true,
		"$(cpu.busy)":                              true,
		"$(cpu.busy)+$(cpu.idle)-$(cpu.nice)":      true,
		"$(cpu.busy)>=80":                          true,
		"($(cpu.busy)+$(cpu.idle)-$(cpu.nice))>80": true,
		"$(qps/module=judge,project=falcon)":       true,
		"($(cpu.idle)+$(cpu.busy))=100":            true,

		// false
		"$((cpu.busy)":                     false,
		"$(cpu.idle)+$(cpu.busy)>40":       false,
		"($(cpu.idle)+$(cpu.busy)-60)>100": false,
	}

	for key, val := range expressionMap {
		if st := expressionValid(key); st != val {
			t.Errorf("func expressionValid() failure")
		}
	}
}
