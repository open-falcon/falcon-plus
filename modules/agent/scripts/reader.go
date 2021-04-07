package scripts

import (
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const (
	ScriptResultTypeJson = "json"
	ScriptResultTypeLine = "line"
	MetricTypeGAUGE = "GAUGE"
	MetricTypeCOUNTER = "COUNTER"
)

// return: dict{./scripts/ntp_60_line.py : *Script}
func ListScripts() map[string]*Script {
	ret := make(map[string]*Script)

	script_path := g.Config().Script.Dir
	fs, err := ioutil.ReadDir(script_path)
	if err != nil {
		log.Println("can not list script files under", script_path)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		filename := f.Name()

		arr := strings.Split(filename, "_")
		larr := len(arr)
		if larr < 3 {
			continue
		}

		// filename should be: $name_$cycle_$resulttype.py
		cycle, err := strconv.Atoi(arr[larr-2])
		if err != nil {
			log.Println(err)
			continue
		}

		resultTypeArr := strings.Split(arr[larr-1], ".")
		if len(resultTypeArr)<1 {
			continue
		}
		resultType := resultTypeArr[0]
		if resultType != ScriptResultTypeJson && resultType!= ScriptResultTypeLine {
			continue
		}

		plugin := &Script{FilePath: filename, MTime: f.ModTime().Unix(), Cycle: cycle,
			Args: "", ResultType:resultType}
		ret[filename] = plugin
	}

	return ret
}
