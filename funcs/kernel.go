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
	}

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
	}

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
	}

	return []*g.MetricValue{
		GaugeValue("kernel.maxfiles", maxFiles),
		GaugeValue("kernel.maxproc", maxProc),
		GaugeValue("kernel.files.allocated", allocateFiles),
		GaugeValue("kernel.files.left", maxFiles-allocateFiles),
	}
}
