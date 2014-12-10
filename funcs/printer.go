package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
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
	}

	err = UpdateDiskStats()
	if err != nil {
		fmt.Println(err)
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

	fmt.Println("all metric collector successfully")
}
