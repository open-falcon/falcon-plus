package funcs

import (
	"os"
	"testing"
)

func TestGetAgentMem(t *testing.T) {
	_ = os.RemoveAll(cgroupRoot)
	// create falcon-agent cgroup dir
	err := os.Mkdir(cgroupRoot, 751)
	if err != nil {
		t.Error(err)
		_ = os.RemoveAll(cgroupRoot)
		return
	}
	_, err = GetAgentMem()
	if err != nil {
		t.Error(err)
		_ = os.RemoveAll(cgroupRoot)
		return
	}
	_ = os.RemoveAll(cgroupRoot)
	return
}
