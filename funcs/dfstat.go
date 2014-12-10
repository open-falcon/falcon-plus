package funcs

import (
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/nux"
	"log"
)

func DeviceMetrics() []*g.MetricValue {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Println(err)
		return nil
	}

	var ret []*g.MetricValue = make([]*g.MetricValue, 0)
	for idx := range mountPoints {
		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2])
		if err != nil {
			log.Println(err)
			continue
		}

		tags := fmt.Sprintf("mount=%s,fstype=%s", du.FsFile, du.FsVfstype)
		ret = append(ret, GaugeValue("df.bytes.total", du.BlocksAll, tags))
		ret = append(ret, GaugeValue("df.bytes.used", du.BlocksUsed, tags))
		ret = append(ret, GaugeValue("df.bytes.free", du.BlocksFree, tags))
		ret = append(ret, GaugeValue("df.bytes.used.percent", du.BlocksUsedPercent, tags))
		ret = append(ret, GaugeValue("df.bytes.free.percent", du.BlocksFreePercent, tags))
		ret = append(ret, GaugeValue("df.inodes.total", du.InodesAll, tags))
		ret = append(ret, GaugeValue("df.inodes.used", du.InodesUsed, tags))
		ret = append(ret, GaugeValue("df.inodes.free", du.InodesFree, tags))
		ret = append(ret, GaugeValue("df.inodes.used.percent", du.InodesUsedPercent, tags))
		ret = append(ret, GaugeValue("df.inodes.free.percent", du.InodesFreePercent, tags))

	}

	return ret
}
