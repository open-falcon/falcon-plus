package g

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func setup() {
	os.MkdirAll("falcon-module/config", 755)
	os.MkdirAll("falcon-module-without-cfg/config", 755)
	os.Create("falcon-module/config/cfg.json")
	cfgOf = map[string]string{
		"falcon-module":             "falcon-module/config/cfg.json",
		"falcon-module-without-cfg": "falcon-module-without-cfg/config/cfg.json",
	}
	Modules = map[string]bool{
		"falcon-module":             true,
		"falcon-module-without-cfg": true,
	}
	AllModulesInOrder = []string{
		"1st-module",
		"2nd-module",
		"3rd-module",
	}
}

func teardown() {
	os.RemoveAll("falcon-module")
	os.RemoveAll("falcon-module-without-cfg")
}

func TestRel(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
		os.Exit(1)
	}
	tests := []struct {
		input    string
		expected string
	}{
		{"falcon-module/config/cfg.json", "falcon-module/config/cfg.json"},
		{"./falcon-module/config/cfg.json", "falcon-module/config/cfg.json"},
		{filepath.Join(wd, "falcon-module/config/cfg.json"), "falcon-module/config/cfg.json"},
	}

	for i, v := range tests {
		actual := Rel(v.input)
		expected := v.expected
		t.Logf("Check case %d: %s(actual) == %s(expected)", i, actual, expected)
		if actual != expected {
			t.Errorf("Error on case %d: %s(actual) != %s(expected)", i, actual, expected)
		}
	}
}

func TestHasCfg(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"falcon-module", true},
		{"falcon-module-without-cfg", false},
	}

	for i, v := range tests {
		actual := HasCfg(v.input)
		expected := v.expected
		t.Logf("Check case %v: %v(actual) == %v(expected)", i, actual, expected)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Error on case %v: %v(actual) != %v(expected)", i, actual, expected)
		}
	}
}

func TestPreqOrder(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"2nd-module", "1st-module"}, []string{"1st-module", "2nd-module"}},
		{[]string{"2nd-module", "1st-module", "3rd-module"}, []string{"1st-module", "2nd-module", "3rd-module"}},
		{[]string{"3rd-module", "2nd-module", "1st-module"}, []string{"1st-module", "2nd-module", "3rd-module"}},
		{[]string{"3rd-module", "other-module", "1st-module", "2nd-module"}, []string{"1st-module", "2nd-module", "3rd-module", "other-module"}},
		{[]string{"other-module", "1st-module", "2nd-module", "3rd-module"}, []string{"1st-module", "2nd-module", "3rd-module", "other-module"}},
	}
	for i, v := range tests {
		actual := PreqOrder(v.input)
		expected := v.expected
		t.Logf("Check case %d: %s(actual) == %s(expected)", i, actual, expected)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Error on case %d: %s(actual) != %s(expected)", i, actual, expected)
		}
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
