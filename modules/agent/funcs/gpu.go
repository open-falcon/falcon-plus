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
	log "github.com/Sirupsen/logrus"

	"github.com/mindprince/gonvml"
	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/modules/agent/g"
)

// 需要load libnvidia-ml.so.1库
func GpuMetrics() (L []*model.MetricValue) {

	if err := gonvml.Initialize(); err != nil {
		debug := g.Config().Debug
		if debug {
			log.Println("Initialize error: ", err)
		}
		return
	}

	defer gonvml.Shutdown()

	count, err := gonvml.DeviceCount()
	if err != nil {
		log.Println("DeviceCount error: ", err)
		return
	}

	if count == 0 {
		return
	}

	temperature := uint(0)
	totalMemory := uint64(0)
	usedMemory := uint64(0)
	gpuUtilization := uint(0)
	memoryUtilization := uint(0)
	powerUsage := uint(0)
	allUtilization := uint(0)
	allMemoryUtilization := uint(0)

	for i := 0; i < int(count); i++ {
		dev, err := gonvml.DeviceHandleByIndex(uint(i))
		if err != nil {
			log.Println("DeviceHandleByIndex error:", err)
			continue
		}

		uuid, err := dev.UUID()
		if err != nil {
			log.Println("dev.UUID error", err)
		}

		tag := "uuid=" + uuid

		// 不是所有gpu都有风扇
		fanSpeed, err := dev.FanSpeed()
		if err != nil {
			log.Println("dev.FanSpeed error: ", err)
		} else {
			L = append(L, GaugeValue("gpu.fan.speed", fanSpeed, tag))
		}

		temperature, err = dev.Temperature()
		if err != nil {
			log.Println("dev.Temperature error: ", err)
			continue
		}

		totalMemory, usedMemory, err = dev.MemoryInfo()
		if err != nil {
			log.Println("dev.MemoryInfo error: ", err)
			continue
		}

		// 单位换算为兆
		totalBillion := float64(totalMemory / 1024 / 1024)
		usedBillion := float64(usedMemory / 1024 / 1024)

		gpuUtilization, memoryUtilization, err = dev.UtilizationRates()
		if err != nil {
			log.Println("dev.UtilizationRates error: ", err)
			continue
		}

		allUtilization += gpuUtilization
		allMemoryUtilization += memoryUtilization

		powerUsage, err = dev.PowerUsage()
		if err != nil {
			log.Println("dev.PowerUsage error: ", err)
		}

		// 单位换算为瓦特
		powerWatt := float64(powerUsage / 1000)

		L = append(L, GaugeValue("gpu.temperature", temperature, tag))
		L = append(L, GaugeValue("gpu.memory.total", totalBillion, tag))
		L = append(L, GaugeValue("gpu.memory.used", usedBillion, tag))
		L = append(L, GaugeValue("gpu.memory.util", memoryUtilization, tag))
		L = append(L, GaugeValue("gpu.util", gpuUtilization, tag))
		L = append(L, GaugeValue("gpu.power.usage", powerWatt, tag))
	}

	L = append(L, GaugeValue("gpu.count", count))
	L = append(L, GaugeValue("gpu.util.avg", allUtilization/count))
	L = append(L, GaugeValue("gpu.memory.util.avg", allMemoryUtilization/count))
	return L
}
