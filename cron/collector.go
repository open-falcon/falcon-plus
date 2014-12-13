package cron

import (
	"github.com/open-falcon/agent/funcs"
	"github.com/open-falcon/agent/g"
	"log"
	"os"
	"time"
)

func InitDataHistory() {
	for {
		funcs.UpdateCpuStat()
		funcs.UpdateDiskStats()
		time.Sleep(time.Second)
	}
}

func Collect() {

	if !g.Config().Transfer.Enabled {
		return
	}

	if g.Config().Transfer.Addr == "" {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*g.MetricValue) {

	for {
	REST:
		time.Sleep(time.Duration(sec) * time.Second)

		hostname, err := os.Hostname()
		if err != nil {
			log.Println("os.Hostname() fail:", err)
			goto REST
		}

		funcIdxItems := make(map[int][]*g.MetricValue, len(fns))

		for idx, fn := range fns {
			items := fn()
			if items == nil {
				continue
			}

			if len(items) == 0 {
				continue
			}

			funcIdxItems[idx] = items
		}

		size := 0
		for _, val := range funcIdxItems {
			size += len(val)
		}

		L := make([]*g.MetricValue, size)

		i := 0
		for _, metrics := range funcIdxItems {
			for _, mv := range metrics {
				L[i] = mv
				i++
			}
		}

		now := time.Now().Unix()
		for j := 0; j < size; j++ {
			L[j].Step = sec
			L[j].Endpoint = hostname
			L[j].Timestamp = now
		}

		g.SendToTransfer(L)

	}
}
