package funcs

import (
	"fmt"
	"github.com/toolkits/nux"
	"github.com/toolkits/sys"
)

func CheckCollector() {

	output := make(map[string]bool)

	_, procStatErr := nux.CurrentProcStat()
	_, listDiskErr := nux.ListDiskStats()
	ports, listeningPortsErr := nux.ListeningPorts()
	procs, psErr := nux.AllProcs()

	_, duErr := sys.CmdOut("du", "--help")

	output["kernel  "] = len(KernelMetrics()) > 0
	output["df.bytes"] = DeviceMetricsCheck()
	output["net.if  "] = len(CoreNetMetrics([]string{})) > 0
	output["loadavg "] = len(LoadAvgMetrics()) > 0
	output["cpustat "] = procStatErr == nil
	output["disk.io "] = listDiskErr == nil
	output["memory  "] = len(MemMetrics()) > 0
	output["netstat "] = len(NetstatMetrics()) > 0
	output["ss -s   "] = len(SocketStatSummaryMetrics()) > 0
	output["ss -tln "] = listeningPortsErr == nil && len(ports) > 0
	output["ps aux  "] = psErr == nil && len(procs) > 0
	output["du -bs  "] = duErr == nil

	for k, v := range output {
		status := "fail"
		if v {
			status = "ok"
		}
		fmt.Println(k, "...", status)
	}
}
