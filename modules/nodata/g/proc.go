package g

import (
	nproc "github.com/toolkits/proc"
	"log"
)

// counter
var (
	ConfigCronCnt = nproc.NewSCounterQps("ConfigCronCnt")
	ConfigLastTs  = nproc.NewSCounterBase("ConfigLastTs")
	ConfigLastCnt = nproc.NewSCounterBase("ConfigLastCnt")

	CollectorCronCnt = nproc.NewSCounterQps("CollectorCronCnt")
	CollectorLastTs  = nproc.NewSCounterBase("CollectorLastTs")
	CollectorLastCnt = nproc.NewSCounterBase("CollectorLastCnt")
	CollectorCnt     = nproc.NewSCounterQps("CollectorCnt")

	JudgeCronCnt = nproc.NewSCounterQps("JudgeCronCnt")
	JudgeLastTs  = nproc.NewSCounterBase("JudgeLastTs")

	SenderCronCnt = nproc.NewSCounterQps("SenderCronCnt")
	SenderLastTs  = nproc.NewSCounterBase("SenderLastTs")
	SenderCnt     = nproc.NewSCounterQps("SenderCnt")
)

// flood
var (
	FloodRate = nproc.NewSCounterBase("FloodRate")
	Threshold = nproc.NewSCounterBase("Threshold")
	Blocking  = nproc.NewSCounterBase("nodata.blocking")
)

func StartProc() {
	log.Println("g.StartProc ok")
}

func GetAllCounters() []interface{} {
	ret := make([]interface{}, 0)

	ret = append(ret, ConfigCronCnt.Get())
	ret = append(ret, ConfigLastTs.Get())
	ret = append(ret, ConfigLastCnt.Get())

	ret = append(ret, CollectorCronCnt.Get())
	ret = append(ret, CollectorLastCnt.Get())
	ret = append(ret, CollectorLastTs.Get())
	ret = append(ret, CollectorCnt.Get())

	ret = append(ret, JudgeCronCnt.Get())
	ret = append(ret, JudgeLastTs.Get())

	ret = append(ret, SenderCronCnt.Get())
	ret = append(ret, SenderLastTs.Get())
	ret = append(ret, SenderCnt.Get())

	ret = append(ret, FloodRate.Get())
	ret = append(ret, Threshold.Get())
	ret = append(ret, Blocking.Get())

	return ret
}
