package nux

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/toolkits/file"
)

type Loadavg struct {
	Avg1min          float64
	Avg5min          float64
	Avg15min         float64
	RunningProcesses int64
	TotalProcesses   int64
}

func (load *Loadavg) String() string {
	return fmt.Sprintf("<1min:%f, 5min:%f, 15min:%f, processes:%d/%d>", load.Avg1min, load.Avg5min, load.Avg15min, load.RunningProcesses, load.TotalProcesses)
}

func LoadAvg() (*Loadavg, error) {

	loadAvg := Loadavg{}

	data, err := file.ToTrimString(Root() + "/proc/loadavg")
	if err != nil {
		return nil, err
	}

	L := strings.Fields(data)
	if loadAvg.Avg1min, err = strconv.ParseFloat(L[0], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg5min, err = strconv.ParseFloat(L[1], 64); err != nil {
		return nil, err
	}
	if loadAvg.Avg15min, err = strconv.ParseFloat(L[2], 64); err != nil {
		return nil, err
	}
	processes := strings.SplitN(L[3], "/", 2)
	if len(processes) != 2 {
		return nil, errors.New("invalid loadavg " + data)
	}
	if loadAvg.RunningProcesses, err = strconv.ParseInt(processes[0], 10, 64); err != nil {
		return nil, err
	}
	if loadAvg.TotalProcesses, err = strconv.ParseInt(processes[1], 10, 64); err != nil {
		return nil, err
	}

	return &loadAvg, nil
}
