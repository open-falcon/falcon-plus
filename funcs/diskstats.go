package funcs

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
	"strings"
	"sync"
)

var (
	diskStatsMap = make(map[string][2]*nux.DiskStats)
	dsLock       = new(sync.RWMutex)
)

func UpdateDiskStats() error {
	dsList, err := nux.ListDiskStats()
	if err != nil {
		return err
	}

	dsLock.Lock()
	defer dsLock.Unlock()
	for i := 0; i < len(dsList); i++ {
		device := dsList[i].Device
		diskStatsMap[device] = [2]*nux.DiskStats{dsList[i], diskStatsMap[device][0]}
	}
	return nil
}

func IOReadRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadRequests - arr[1].ReadRequests
}

func IOReadMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadMerged - arr[1].ReadMerged
}

func IOReadSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].ReadSectors - arr[1].ReadSectors
}

func IOMsecRead(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecRead - arr[1].MsecRead
}

func IOWriteRequests(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteRequests - arr[1].WriteRequests
}

func IOWriteMerged(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteMerged - arr[1].WriteMerged
}

func IOWriteSectors(arr [2]*nux.DiskStats) uint64 {
	return arr[0].WriteSectors - arr[1].WriteSectors
}

func IOMsecWrite(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWrite - arr[1].MsecWrite
}

func IOMsecTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecTotal - arr[1].MsecTotal
}

func IOMsecWeightedTotal(arr [2]*nux.DiskStats) uint64 {
	return arr[0].MsecWeightedTotal - arr[1].MsecWeightedTotal
}

func IODelta(device string, f func([2]*nux.DiskStats) uint64) uint64 {
	val, ok := diskStatsMap[device]
	if !ok {
		return 0
	}

	if val[1] == nil {
		return 0
	}
	return f(val)
}

func DiskIOMetrics() []*g.MetricValue {

	ret := make([]*g.MetricValue, 0)

	dsList, err := nux.ListDiskStats()
	if err != nil {
		log.Println(err)
		return ret
	}

	for _, ds := range dsList {
		if !shouldHandle(ds.Device) {
			continue
		}

		device := "device=" + ds.Device

		ret = append(ret, CounterValue("disk.io.read_requests", ds.ReadRequests, device))
		ret = append(ret, CounterValue("disk.io.read_merged", ds.ReadMerged, device))
		ret = append(ret, CounterValue("disk.io.read_sectors", ds.ReadSectors, device))
		ret = append(ret, CounterValue("disk.io.msec_read", ds.MsecRead, device))
		ret = append(ret, CounterValue("disk.io.write_requests", ds.WriteRequests, device))
		ret = append(ret, CounterValue("disk.io.write_merged", ds.WriteMerged, device))
		ret = append(ret, CounterValue("disk.io.write_sectors", ds.WriteSectors, device))
		ret = append(ret, CounterValue("disk.io.msec_write", ds.MsecWrite, device))
		ret = append(ret, CounterValue("disk.io.ios_in_progress", ds.IosInProgress, device))
		ret = append(ret, CounterValue("disk.io.msec_total", ds.MsecTotal, device))
		ret = append(ret, CounterValue("disk.io.msec_weighted_total", ds.MsecWeightedTotal, device))
	}
	return ret
}

func IOStatsMetrics() []*g.MetricValue {
	ret := make([]*g.MetricValue, 0)

	dsLock.RLock()
	defer dsLock.RUnlock()

	for device, _ := range diskStatsMap {
		if !shouldHandle(device) {
			continue
		}

		tags := "device=" + device
		rio := IODelta(device, IOReadRequests)
		wio := IODelta(device, IOWriteRequests)
		delta_rsec := IODelta(device, IOReadSectors)
		delta_wsec := IODelta(device, IOWriteSectors)
		ruse := IODelta(device, IOMsecRead)
		wuse := IODelta(device, IOMsecWrite)
		use := IODelta(device, IOMsecTotal)
		n_io := rio + wio
		avgrq_sz := 0.0
		await := 0.0
		svctm := 0.0
		if n_io != 0 {
			avgrq_sz = float64(delta_rsec+delta_wsec) / float64(n_io)
			await = float64(ruse+wuse) / float64(n_io)
			svctm = float64(use) / float64(n_io)
		}

		ret = append(ret, GaugeValue("disk.io.read_bytes", float64(delta_rsec)*512.0, tags))
		ret = append(ret, GaugeValue("disk.io.write_bytes", float64(delta_wsec)*512.0, tags))
		ret = append(ret, GaugeValue("disk.io.avgrq_sz", avgrq_sz, tags))
		ret = append(ret, GaugeValue("disk.io.avgqu-sz", float64(IODelta(device, IOMsecWeightedTotal))/1000.0, tags))
		ret = append(ret, GaugeValue("disk.io.await", await, tags))
		ret = append(ret, GaugeValue("disk.io.svctm", svctm, tags))
		ret = append(ret, GaugeValue("disk.io.util", float64(use)/10.0, tags))
	}

	return ret
}

func shouldHandle(device string) bool {
	normal := len(device) == 3 && (strings.HasPrefix(device, "sd") || strings.HasPrefix(device, "vd"))
	aws := len(device) == 4 && strings.HasPrefix(device, "xvd")
	return normal || aws
}
