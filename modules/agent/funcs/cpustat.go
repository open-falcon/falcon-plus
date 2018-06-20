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
	"sync"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/toolkits/nux"
)

const (
	historyCount int = 2
)

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}

	psLock.Lock()
	defer psLock.Unlock()
	for i := historyCount - 1; i > 0; i-- {
		procStatHistory[i] = procStatHistory[i-1]
	}

	procStatHistory[0] = ps
	return nil
}

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func historyPrepared() bool {
	return procStatHistory[1] != nil
}

func CpuUsagesSummary() (cpuUsages []float64, switches uint64, prepared bool) {
	psLock.RLock()
	defer psLock.RUnlock()

	cpuUsages = make([]float64, 0, 10)
	switches = 0
	prepared = historyPrepared()

	if !prepared {
		return
	}

	dt := deltaTotal()

	if dt == 0 {
		idle := 0.0
		busy := 100.0 - idle
		user := 0.0
		nice := 0.0
		system := 0.0
		iowait := 0.0
		irq := 0.0
		softirq := 0.0
		steal := 0.0
		guest := 0.0

		cpuUsages = append(cpuUsages, idle, busy, user, nice, system, iowait, irq, softirq, steal, guest)
		switches = procStatHistory[0].Ctxt
	} else {
		invQuotient := 100.00 / float64(dt)

		idle := float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
		busy := 100.0 - idle
		user := float64(procStatHistory[0].Cpu.User-procStatHistory[1].Cpu.User) * invQuotient
		nice := float64(procStatHistory[0].Cpu.Nice-procStatHistory[1].Cpu.Nice) * invQuotient
		system := float64(procStatHistory[0].Cpu.System-procStatHistory[1].Cpu.System) * invQuotient
		iowait := float64(procStatHistory[0].Cpu.Iowait-procStatHistory[1].Cpu.Iowait) * invQuotient
		irq := float64(procStatHistory[0].Cpu.Irq-procStatHistory[1].Cpu.Irq) * invQuotient
		softirq := float64(procStatHistory[0].Cpu.SoftIrq-procStatHistory[1].Cpu.SoftIrq) * invQuotient
		steal := float64(procStatHistory[0].Cpu.Steal-procStatHistory[1].Cpu.Steal) * invQuotient
		guest := float64(procStatHistory[0].Cpu.Guest-procStatHistory[1].Cpu.Guest) * invQuotient

		cpuUsages = append(cpuUsages, idle, busy, user, nice, system, iowait, irq, softirq, steal, guest)
		switches = procStatHistory[0].Ctxt
	}

	return
}

func CpuMetrics() []*model.MetricValue {
	cpuUsages, currentCpuSwitches, prepared := CpuUsagesSummary()

	if !prepared {
		return []*model.MetricValue{}
	}

	idle := GaugeValue("cpu.idle", cpuUsages[0])
	busy := GaugeValue("cpu.busy", cpuUsages[1])
	user := GaugeValue("cpu.user", cpuUsages[2])
	nice := GaugeValue("cpu.nice", cpuUsages[3])
	system := GaugeValue("cpu.system", cpuUsages[4])
	iowait := GaugeValue("cpu.iowait", cpuUsages[5])
	irq := GaugeValue("cpu.irq", cpuUsages[6])
	softirq := GaugeValue("cpu.softirq", cpuUsages[7])
	steal := GaugeValue("cpu.steal", cpuUsages[8])
	guest := GaugeValue("cpu.guest", cpuUsages[9])
	switches := CounterValue("cpu.switches", currentCpuSwitches)
	return []*model.MetricValue{idle, busy, user, nice, system, iowait, irq, softirq, steal, guest, switches}
}
