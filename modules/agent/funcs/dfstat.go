// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package funcs

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
	"github.com/toolkits/nux"
)

func DeviceMetrics() (L []*model.MetricValue) {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Error("collect device metrics fail:", err)
		return
	}

	var myMountPoints map[string]bool = make(map[string]bool)

	if len(g.Config().Collector.MountPoint) > 0 {
		for _, mp := range g.Config().Collector.MountPoint {
			myMountPoints[mp] = true
		}
	}

	var diskTotal uint64 = 0
	var diskUsed uint64 = 0

	for idx := range mountPoints {
		fsSpec, fsFile, fsVfstype := mountPoints[idx][0], mountPoints[idx][1], mountPoints[idx][2]
		if len(myMountPoints) > 0 {
			if _, ok := myMountPoints[fsFile]; !ok {
				log.Debug("mount point not matched with config", fsFile, "ignored.")
				continue
			}
		}

		var du *nux.DeviceUsage
		du, err = nux.BuildDeviceUsage(fsSpec, fsFile, fsVfstype)
		if err != nil {
			log.Error(err)
			continue
		}

		if du.BlocksAll == 0 {
			continue
		}

		if du.FsVfstype == "iso9660" {
			continue
		}

		diskTotal += du.BlocksAll
		diskUsed += du.BlocksUsed

		tags := fmt.Sprintf("mount=%s,fstype=%s", du.FsFile, du.FsVfstype)
		L = append(L, GaugeValue("df.bytes.total", du.BlocksAll, tags))
		L = append(L, GaugeValue("df.bytes.used", du.BlocksUsed, tags))
		L = append(L, GaugeValue("df.bytes.free", du.BlocksFree, tags))
		L = append(L, GaugeValue("df.bytes.used.percent", du.BlocksUsedPercent, tags))
		L = append(L, GaugeValue("df.bytes.free.percent", du.BlocksFreePercent, tags))

		if du.InodesAll == 0 {
			continue
		}

		L = append(L, GaugeValue("df.inodes.total", du.InodesAll, tags))
		L = append(L, GaugeValue("df.inodes.used", du.InodesUsed, tags))
		L = append(L, GaugeValue("df.inodes.free", du.InodesFree, tags))
		L = append(L, GaugeValue("df.inodes.used.percent", du.InodesUsedPercent, tags))
		L = append(L, GaugeValue("df.inodes.free.percent", du.InodesFreePercent, tags))

	}

	if len(L) > 0 && diskTotal > 0 {
		L = append(L, GaugeValue("df.statistics.total", float64(diskTotal)))
		L = append(L, GaugeValue("df.statistics.used", float64(diskUsed)))
		L = append(L, GaugeValue("df.statistics.used.percent", float64(diskUsed)*100.0/float64(diskTotal)))
	}

	return
}

func DeviceMetricsCheck() bool {
	mountPoints, err := nux.ListMountPoint()

	if err != nil {
		log.Error("collect device metrics fail:", err)
		return false
	}

	if len(mountPoints) <= 0 {
		return false
	}

	return true
}
