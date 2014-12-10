package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func KernelMetrics() []*g.MetricValue {

	ret := []*g.MetricValue{}

	maxFiles, err := nux.KernelMaxFiles()
	if err != nil {
		log.Println(err)
		return ret
	}

	ret = append(ret, GaugeValue("kernel.maxfiles", maxFiles))

	maxProc, err := nux.KernelMaxProc()
	if err != nil {
		log.Println(err)
		return ret
	}

	ret = append(ret, GaugeValue("kernel.maxproc", maxProc))

	allocateFiles, err := nux.KernelAllocateFiles()
	if err != nil {
		log.Println(err)
		return ret
	}

	ret = append(ret, GaugeValue("kernel.files.allocated", allocateFiles))
	ret = append(ret, GaugeValue("kernel.files.left", maxFiles-allocateFiles))
	return ret
}
