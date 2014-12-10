package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"os"
	"time"
)

func p(item *g.MetricValue) {
	fmt.Printf("%s=%v[tags:%s]\n", item.Metric, item.Value, item.Tags)
}

func PrintAll() {
	for _, item := range AgentMetrics() {
		p(item)
	}

	for _, item := range KernelMetrics() {
		p(item)
	}

	for _, item := range DeviceMetrics() {
		p(item)
	}

	for _, item := range CoreNetMetrics([]string{}) {
		p(item)
	}

	for _, item := range LoadAvgMetrics() {
		p(item)
	}

	err := UpdateCpuStat()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = UpdateDiskStats()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	time.Sleep(time.Second)

	UpdateCpuStat()
	UpdateDiskStats()

	for _, item := range CpuMetrics() {
		p(item)
	}

	for _, item := range DiskIOMetrics() {
		p(item)
	}

	for _, item := range IOStatsMetrics() {
		p(item)
	}

	for _, item := range MemMetrics() {
		p(item)
	}

	for _, item := range NetstatMetrics() {
		p(item)
	}

	for _, item := range SocketStatSummaryMetrics() {
		p(item)
	}

	fmt.Print("Listening ports: ")
	fmt.Println(nux.ListeningPorts())

	procs, err := nux.AllProcs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cnt := 0
	for _, item := range procs {
		fmt.Println(item)
		cnt++
		if cnt == 10 {
			fmt.Println("...")
			break
		}
	}

	fmt.Println("all metric collector successfully")
}
