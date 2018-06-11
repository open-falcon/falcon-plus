package nux

import (
	"bufio"
	"bytes"
	"strconv"
	"strings"

	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

func SocketStatSummary() (m map[string]uint64, err error) {
	m = make(map[string]uint64)
	var bs []byte
	bs, err = sys.CmdOutBytes("sh", "-c", "ss -s")
	if err != nil {
		return
	}

	reader := bufio.NewReader(bytes.NewBuffer(bs))

	// ignore the first line
	line, e := file.ReadLine(reader)
	if e != nil {
		return m, e
	}

	for {
		line, err = file.ReadLine(reader)
		if err != nil {
			return
		}

		lineStr := string(line)
		if strings.HasPrefix(lineStr, "TCP") {
			left := strings.Index(lineStr, "(")
			right := strings.Index(lineStr, ")")
			if left < 0 || right < 0 {
				continue
			}

			content := lineStr[left+1 : right]
			arr := strings.Split(content, ", ")
			for _, val := range arr {
				fields := strings.Fields(val)
				if fields[0] == "timewait" {
					timewait_arr := strings.Split(fields[1], "/")
					m["timewait"], _ = strconv.ParseUint(timewait_arr[0], 10, 64)
					m["slabinfo.timewait"], _ = strconv.ParseUint(timewait_arr[1], 10, 64)
					continue
				}
				m[fields[0]], _ = strconv.ParseUint(fields[1], 10, 64)
			}
			return
		}
	}

	return
}
