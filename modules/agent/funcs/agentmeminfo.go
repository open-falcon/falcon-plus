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

type Mem struct {
	Buffers   uint64
	Cached    uint64
	MemTotal  uint64
	MemFree   uint64
	SwapTotal uint64
	SwapUsed  uint64
	SwapFree  uint64
	VmRSS     uint64
}

var Multi uint64 = 1024

var WANT = map[string]struct{}{
	"Buffers":   struct{}{},
	"Cached":    struct{}{},
	"MemTotal":  struct{}{},
	"MemFree":   struct{}{},
	"SwapTotal": struct{}{},
	"SwapUsed":  struct{}{},
	"SwapFree":  struct{}{},
	"VmRSS":     struct{}{},
}

// ProcessInfo 定义进程信息
type ProcessInfo struct {
	name  string
	pid   int
	ppid  int
	state string
}

// ProcessManager
type ProcessManager struct {
	handler *os.File
}

func AgentMemInfo() (*Mem, error) {
	proc := ProcessManager{}
	pid, err := proc.CheckProc()
	if err != nil {
		return nil, err
	}
	var strs = strconv.Itoa(pid)
	contents, err := ioutil.ReadFile("/proc" + strs + "status")
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	memInfo := &Mem{}
	reader := bufio.NewReader(bytes.NewBuffer(contents))

	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}
		fields := strings.Fields(string(line))
		fieldName := fields[0]
		_, ok := WANT[fieldName]
		if ok && len(fields) == 3 {
			val, numerr := strconv.ParseUint(fields[1], 10, 64)
			if numerr != nil {
				continue
			}
			switch fieldName {
			case "VmRSS":
				memInfo.VmRSS = val / Multi
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return memInfo, nil
}

func (proc *ProcessManager) CheckProc() (int, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return 0, err
	}
	defer func() {
		d.Close()
	}()
	proc.handler = d
	for {
		fileList, err := proc.handler.Readdir(10)
		if err != nil && err != io.EOF {
			return 0, err
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
			p := ProcessInfo{pid: int(pid)}
			if err := p.Load(); err == nil {
				return p.pid, nil
			} else {
				continue
			}
		}
	}
}

// 确认下
func (p *ProcessInfo) Load() error {
	dataBytes, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", p.pid))
	if err != nil {
		return err
	}
	data := string(dataBytes)
	start := strings.IndexRune(data, '(') + 1
	end := strings.IndexRune(data[start:], ')')
	p.name = data[start : start+end]
	if p.name == "falcon-agent" {
		result := strings.Split(string(data[start+end+2:]), " ")
		if len(result) < 2 {
			return errors.New("length not right")
		}
		p.state = result[0]
		if ppid, err := strconv.Atoi(result[2]); err == nil {
			p.ppid = ppid
		}
	}
	return errors.New("the proc not exist")
}
