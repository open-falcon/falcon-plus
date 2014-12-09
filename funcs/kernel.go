package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func KernelMetrics() []*g.MetricValue {

	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Println(err)
		return []*g.MetricValue{}
	}

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
		return []*g.MetricValue{}
	}

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
		return []*g.MetricValue{}
	}

	return []*g.MetricValue{
		GaugeValue("kernel.maxfiles", maxFiles),
		GaugeValue("kernel.maxproc", maxProc),
		GaugeValue("kernel.files.allocated", allocateFiles),
		GaugeValue("kernel.files.left", maxFiles-allocateFiles),
	}
}
