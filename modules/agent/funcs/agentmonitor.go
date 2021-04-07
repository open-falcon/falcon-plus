package funcs

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/open-falcon/falcon-plus/g"
	"github.com/toolkits/file"
)

var (
	memLimit   int    = 200 // agent最大可以使用的内存，单位MB
	cgroupRoot string = "/sys/fs/cgroup/memory/falcon-agent"
)

const (
	procsFile = "cgroup.procs"
	memStat   = "memory.stat"
	mb        = 1024 * 1024
)

// InitCgroup init falcon-agent
func InitCgroup() {
	pid := g.Pid("agent")
	_ = os.RemoveAll(cgroupRoot)
	// create falcon-agent cgroup dir
	err := os.Mkdir(cgroupRoot, 751)
	if err != nil {
		fmt.Println("falcon-agent cgroup init failed", err)
		return
	}
	// set memory limit
	pPath := filepath.Join(cgroupRoot, procsFile)
	err = ioutil.WriteFile(pPath, []byte(fmt.Sprintf("%s", pid)), 644)
	if err != nil {
		fmt.Println("falcon-agent cgroup write cgroup.procs failed", err)
		return
	}
}

// GetAgentMem get agent memory info
func GetAgentMem() (int64, error) {
	filePath := filepath.Join(cgroupRoot, memStat)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error: %v", err)
		return 0, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	var agentRSS int64
	for {
		info, err := file.ReadLine(reader)
		if err != nil {
			return agentRSS, err
		}
		fields := strings.Fields(string(info))
		if len(fields) < 2 || fields[0] != "rss" {
			continue
		}
		val, numErr := strconv.ParseInt(fields[1], 10, 64)
		if numErr != nil {
			continue
		}
		agentRSS = val / mb
		break
	}
	return agentRSS, nil
}
