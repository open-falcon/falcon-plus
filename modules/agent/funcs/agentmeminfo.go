package funcs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/toolkits/file"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const Multi uint64 = 1024

// ProcessInfo 定义进程信息
type ProcessInfo struct {
	name  string
	pid   int
	ppid  int
	state string
}

// agentInfo falcon-agent进程信息
var agentInfo *ProcessInfo

func init() {
	agentInfo = getProcInfo("falcon-agent")
	if agentInfo == nil {
		log.Println("not exist falcon-agent proc")
	}
}

func AgentMemInfo() (uint64, error) {
	if agentInfo == nil {
		agentInfo = getProcInfo("falcon-agent")
		if agentInfo==nil {
			return 0, errors.New("not exist falcon-agent proc")
		}
	}
	pid := agentInfo.pid
	contents, err := ioutil.ReadFile("/proc" + strconv.Itoa(pid) + "status")
	if err != nil {
		log.Printf("error: %v", err)
		return 0, err
	}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	var agentVmRSS uint64
	for {
		info, err := file.ReadLine(reader)
		if err != nil {
			break
		}
		fields := strings.Fields(string(info))
		if len(fields) < 2 || fields[0] != "VmRSS" {
			continue
		}
		val, numErr := strconv.ParseUint(fields[1], 10, 64)
		if numErr != nil {
			continue
		}
		agentVmRSS = val / Multi
	}
	if err != io.EOF {
		return 0, err
	}
	return agentVmRSS, nil
}

func getProcInfo(procName string) *ProcessInfo {
	dirs, err := os.Open("/proc")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer dirs.Close()
	for {
		fileList, err := dirs.Readdir(10)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, fi := range fileList {
			if !fi.IsDir() {
				continue
			}
			name := fi.Name()
			if ok, err := regexp.MatchString(`^[0-9]+$`, name); err != nil || !ok {
				continue
			}
			pid, err := strconv.ParseInt(name, 10, 0)
			if err != nil {
				continue
			}
			procData, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
			if err != nil {
				continue
			}
			procInfo := matchProc(procData, procName)
			if procInfo == nil {
				continue
			}
			procInfo.pid=int(pid)
			return procInfo
		}
	}
}

// matchProc 获得相应进程名称的pid的信息
func matchProc(procData []byte, procName string) *ProcessInfo {
	var p = &ProcessInfo{
		name:procName,
	}
	data := string(procData)
	start := strings.IndexRune(data, '(') + 1
	end := strings.IndexRune(data[start:], ')')
	otherInfo := strings.Split(data[start+end+2:], " ")
	if data[start:start+end] != procName || len(otherInfo) < 3 {
		return nil
	}
	p.state = otherInfo[0]
	ppid, err := strconv.Atoi(otherInfo[2])
	if err == nil {
		p.ppid = ppid
	}
	return p
}
