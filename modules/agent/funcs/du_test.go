package funcs

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"os"
	"path/filepath"
	"testing"
)

var testPath = ""

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	testPath = dir

}

func TestDuMetrics(t *testing.T) {
	paths := []string{testPath}
	g.SetDuPaths(paths)
	duMetrics := DuMetrics()
	if duMetrics[0].Value == -1 {
		t.Error("expect metric value > 0 ,but -1")
	}
}

func TestDuMetricsPathNotExist(t *testing.T) {
	paths := []string{"'path does not exist'"}
	g.SetDuPaths(paths)
	duMetrics := DuMetrics()
	if len(duMetrics) != 1 || duMetrics[0].Value != -1 {
		t.Errorf("metric value expect -1, but metric : %v", duMetrics)
	}
}

func TestDuMetricsTimeout(t *testing.T) {
	timeout = 0
	paths := []string{testPath}
	g.SetDuPaths(paths)
	duMetrics := DuMetrics()
	if len(duMetrics) != 1 || duMetrics[0].Value != -1 {
		t.Errorf("metric value expect -1, but metric : %v", duMetrics)
	}
}
