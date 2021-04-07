package scripts

import (
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"strconv"
	"fmt"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

type ScriptScheduler struct {
	Ticker *time.Ticker
	Script *Script
	Quit   chan struct{}
}

func NewScriptScheduler(p *Script) *ScriptScheduler {
	scheduler := ScriptScheduler{Script: p}
	scheduler.Ticker = time.NewTicker(time.Duration(p.Cycle) * time.Second)
	scheduler.Quit = make(chan struct{})
	return &scheduler
}

func (this *ScriptScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				ScriptRun(this.Script)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *ScriptScheduler) Stop() {
	close(this.Quit)
}

// using ',' as the seprator of args and '\,' to espace
func ScriptArgsParse(raw_args string) []string {
	ss := strings.Split(raw_args, "\\,")

	out := [][]string{}
	for _, s := range ss {
		clean_args := []string{}
		for _, arg := range strings.Split(s, ",") {
			arg = strings.Trim(arg, " ")
			arg = strings.Trim(arg, "\"")
			arg = strings.Trim(arg, "'")
			clean_args = append(clean_args, arg)
		}
		out = append(out, clean_args)
	}

	ret := []string{}
	tail := ""

	for _, x := range out {
		for j, y := range x {
			if j == 0 {
				if tail != "" {
					ret = append(ret, tail+","+y)
					tail = ""
				} else {
					ret = append(ret, y)
				}
			} else if j == len(x)-1 {
				tail = y
			} else {
				ret = append(ret, y)
			}
		}
	}

	if tail != "" {
		ret = append(ret, tail)
	}

	return ret
}

func ScriptRun(plugin *Script) {
	timeout := plugin.Cycle*1000 - 500
	fpath := filepath.Join(g.Config().Script.Dir, plugin.FilePath)
	args := plugin.Args

	if !file.IsExist(fpath) {
		log.Printf("no such script: %s(%s)", fpath, args)
		return
	}

	debug := g.Config().Debug
	if debug {
		log.Printf("%s(%s) running...", fpath, args)
	}

	var cmd *exec.Cmd
	if args == "" {
		cmd = exec.Command(fpath)
	} else {
		arg_list := ScriptArgsParse(args)
		cmd = exec.Command(fpath, arg_list...)
	}
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := cmd.Start()
	if err != nil {
		log.Printf("[ERROR] script start fail: %s(%s) , error: %s\n", fpath, args, err)
		return
	}
	if debug {
		log.Printf("script started: %s(%s)", fpath, args)
	}

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)

	errStr := stderr.String()
	if errStr != "" {
		logFile := filepath.Join(g.Config().Script.LogDir, plugin.FilePath+"("+plugin.Args+")"+".stderr.log")
		if _, err = file.WriteString(logFile, errStr); err != nil {
			log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
		}
	}

	if isTimeout {
		// has be killed
		if err == nil && debug {
			log.Println("[INFO] timeout and kill process ", fpath, "(", args, ")", " successfully")
		}

		if err != nil {
			log.Println("[ERROR] kill process ", fpath, "(", args, ")", " occur error:", err)
		}

		return
	}

	if err != nil {
		log.Println("[ERROR] exec script", fpath, "(", args, ")", "fail. error:", err)
		return
	}

	// exec successfully
	data := stdout.Bytes()
	if len(data) == 0 {
		if debug {
			log.Println("[DEBUG] stdout of", fpath, "(", args, ")", "is blank")
		}
		return
	}

	var metrics []*model.MetricValue
	metrics, err = unmarshalScriptResult(data, plugin)
	if err != nil {
		log.Printf("[ERROR] unmarshal stdout of %s(%s) fail. error:%s stdout: \n%s\n", fpath, args, err, stdout.String())
		return
	}

	if len(metrics)!=0 {
		g.SendToTransfer(metrics)	
	}
}

func unmarshalScriptResult(rbytes []byte, plugin *Script) ([]*model.MetricValue, error) {
	var metrics []*model.MetricValue

	// json
	if plugin.ResultType == ScriptResultTypeJson {
		err := json.Unmarshal(rbytes, &metrics)
		return metrics, err
	}

	// line举例:
	// codis_proxy_online: 1 : GAUGE
	// codis_proxy_closed: 0 : COUNTER
	// codis_proxy_ops_total: 6678532816
	// codis_proxy_ops_fails: 520344
	// codis_proxy_ops_qps: 226
	// codis_proxy_ops_redis_errors: 843
	// codis_proxy_sessions_total: 8028100437
	// codis_proxy_sessions_alive: 44
	if plugin.ResultType == ScriptResultTypeLine {
		hostname, err := g.Hostname()
		if err != nil {
			return nil, err
		}
		ts := time.Now().Unix()

		lines := strings.Split(string(rbytes), "\n")
		for i:=0; i<len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			pairs := strings.Split(line, ":")
			if len(pairs) < 2 {
				continue
			}
			metric := strings.TrimSpace(pairs[0])
			valstr := strings.TrimSpace(pairs[1])
			value, err := strconv.ParseFloat(valstr, 64)
			if err != nil {
				continue
			}
			mtype := "GAUGE"
			if len(pairs)>=3 {
				mtypeTmp := strings.TrimSpace(pairs[2])
				if mtypeTmp == MetricTypeCOUNTER {
					mtype = MetricTypeCOUNTER
				}
			}
			metrics = append(metrics, &model.MetricValue{
					Endpoint: hostname,
					Metric: metric,
					Value: value,
					Step: int64(plugin.Cycle),
					Type: mtype,
					Tags: "",
					Timestamp: ts,
				})
		}

		return metrics, nil
	}

	return nil, fmt.Errorf("bad script result type %s", plugin.ResultType)
}

