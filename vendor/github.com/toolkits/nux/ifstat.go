package nux

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/toolkits/file"
	"github.com/toolkits/sys"
)

const (
	BITS_PER_BYTE = 8
	MILLION_BIT   = 1000000
)

type NetIf struct {
	Iface          string
	InBytes        int64
	InPackages     int64
	InErrors       int64
	InDropped      int64
	InFifoErrs     int64
	InFrameErrs    int64
	InCompressed   int64
	InMulticast    int64
	OutBytes       int64
	OutPackages    int64
	OutErrors      int64
	OutDropped     int64
	OutFifoErrs    int64
	OutCollisions  int64
	OutCarrierErrs int64
	OutCompressed  int64
	TotalBytes     int64
	TotalPackages  int64
	TotalErrors    int64
	TotalDropped   int64
	SpeedBits      int64
	InPercent      float64
	OutPercent     float64
}

func (this *NetIf) String() string {
	return fmt.Sprintf("<Iface:%s,InBytes:%d,InPackages:%d...>", this.Iface, this.InBytes, this.InPackages)
}

/*
Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
  eth0: 1990350    2838    0    0    0     0          0         0   401351    2218    0    0    0     0       0          0
    lo:   26105     286    0    0    0     0          0         0    26105     286    0    0    0     0       0          0
*/
func NetIfs(onlyPrefix []string) ([]*NetIf, error) {
	contents, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	ret := []*NetIf{}

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		lineBytes, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		line := string(lineBytes)
		idx := strings.Index(line, ":")
		if idx < 0 {
			continue
		}

		netIf := NetIf{}

		eth := strings.TrimSpace(line[0:idx])
		if len(onlyPrefix) > 0 {
			found := false
			for _, prefix := range onlyPrefix {
				if strings.HasPrefix(eth, prefix) {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		netIf.Iface = eth

		fields := strings.Fields(line[idx+1:])

		if len(fields) != 16 {
			continue
		}

		netIf.InBytes, _ = strconv.ParseInt(fields[0], 10, 64)
		netIf.InPackages, _ = strconv.ParseInt(fields[1], 10, 64)
		netIf.InErrors, _ = strconv.ParseInt(fields[2], 10, 64)
		netIf.InDropped, _ = strconv.ParseInt(fields[3], 10, 64)
		netIf.InFifoErrs, _ = strconv.ParseInt(fields[4], 10, 64)
		netIf.InFrameErrs, _ = strconv.ParseInt(fields[5], 10, 64)
		netIf.InCompressed, _ = strconv.ParseInt(fields[6], 10, 64)
		netIf.InMulticast, _ = strconv.ParseInt(fields[7], 10, 64)

		netIf.OutBytes, _ = strconv.ParseInt(fields[8], 10, 64)
		netIf.OutPackages, _ = strconv.ParseInt(fields[9], 10, 64)
		netIf.OutErrors, _ = strconv.ParseInt(fields[10], 10, 64)
		netIf.OutDropped, _ = strconv.ParseInt(fields[11], 10, 64)
		netIf.OutFifoErrs, _ = strconv.ParseInt(fields[12], 10, 64)
		netIf.OutCollisions, _ = strconv.ParseInt(fields[13], 10, 64)
		netIf.OutCarrierErrs, _ = strconv.ParseInt(fields[14], 10, 64)
		netIf.OutCompressed, _ = strconv.ParseInt(fields[15], 10, 64)

		netIf.TotalBytes = netIf.InBytes + netIf.OutBytes
		netIf.TotalPackages = netIf.InPackages + netIf.OutPackages
		netIf.TotalErrors = netIf.InErrors + netIf.OutErrors
		netIf.TotalDropped = netIf.InDropped + netIf.OutDropped

		speedFile := fmt.Sprintf("/sys/class/net/%s/speed", netIf.Iface)
		if content, err := ioutil.ReadFile(speedFile); err == nil {
			var speed int64
			speed, err = strconv.ParseInt(strings.TrimSpace(string(content)), 10, 64)
			if err != nil {
				netIf.SpeedBits = int64(0)
				netIf.InPercent = float64(0)
				netIf.OutPercent = float64(0)
			} else if speed == 0 {
				netIf.SpeedBits = int64(0)
				netIf.InPercent = float64(0)
				netIf.OutPercent = float64(0)
			} else {
				netIf.SpeedBits = speed * MILLION_BIT
				netIf.InPercent = float64(netIf.InBytes*BITS_PER_BYTE) * 100.0 / float64(netIf.SpeedBits)
				netIf.OutPercent = float64(netIf.OutBytes*BITS_PER_BYTE) * 100.0 / float64(netIf.SpeedBits)
			}
		} else {
			if content, err := sys.CmdOutBytes("ethtool", netIf.Iface); err == nil {
				var speed int64
				var speedStr string

				contentReader := bufio.NewReader(bytes.NewBuffer(content))
				for {
					line, err := file.ReadLine(contentReader)

					if err == io.EOF {
						err = nil
						break
					}

					if err != nil {
						break
					}

					line = bytes.Trim(line, "\t")

					if bytes.HasPrefix(line, []byte("Speed:")) && bytes.HasSuffix(line, []byte("Mb/s")) {
						speedStr = string(line[7 : len(line)-4])
						break
					}
				}

				speed, err = strconv.ParseInt(strings.TrimSpace(speedStr), 10, 64)
				if speedStr == "" || err != nil || speed == 0 {
					netIf.SpeedBits = int64(0)
					netIf.InPercent = float64(0)
					netIf.OutPercent = float64(0)
				} else {
					netIf.SpeedBits = speed * MILLION_BIT
					netIf.InPercent = float64(netIf.InBytes*BITS_PER_BYTE) * 100.0 / float64(netIf.SpeedBits)
					netIf.OutPercent = float64(netIf.OutBytes*BITS_PER_BYTE) * 100.0 / float64(netIf.SpeedBits)
				}
			} else {
				netIf.SpeedBits = int64(0)
				netIf.InPercent = float64(0)
				netIf.OutPercent = float64(0)
			}
		}

		ret = append(ret, &netIf)
	}

	return ret, nil
}
