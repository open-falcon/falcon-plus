package fake_data_gen

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/open-falcon/falcon-plus/modules/api/app/model/alarm"
)

type BossObj struct {
	Platform string
	Province string
	Isp      string
	Idc      string
	Ip       string
	Hostname string
}

var (
	metric   = []string{"cpu.idle", "cpu.busy", "cpu.guest", "cpu.iowait"}
	sfunc    = []string{"all(#3)", "all(#1)"}
	cond     = []string{"3 == 3", "1 >= 0", "0 < 1"}
	note     = []string{"测试1", "note2", "note3"}
	maxStep  = []int{3, 2, 1}
	priority = []int{0, 1, 2, 3, 4}
	status   = []string{"OK", "PROBLEM"}
)

func getMetric() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(metric) - 1)
	return metric[indx]
}

func getFunc() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(getsfunc) - 1)
	return getsfunc[indx]
}

func getCond() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(cond) - 1)
	return cond[indx]
}

func getPriority() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(priority) - 1)
	return note[priority]
}

func getNote() string {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(note) - 1)
	return note[indx]
}

func maxStep() (int, int) {
	rand.Seed(time.Now().UnixNano())
	indx := rand.Intn(len(maxStep) - 1)
	return maxStep[indx], randIntn(maxStep[indx])
}

func main() {
	maxStepn, currentStep := maxStep()
	s := alarm.EventCases{
		ID:          randomdata.StringNumber(12, "-"),
		Endpoint:    fmt.Sprintf("%v_%s", randomdata.StringNumber(3, "-"), randomdata.IpV4Address()),
		Metric:      getMetric(),
		Func:        getFunc(),
		Cond:        getCond(),
		Note:        getNote(),
		MaxStep:     maxStepn,
		CurrentStep: currentStep,
	}
	f, err := os.Create("../fakeData.json")
	defer f.Close()
	nsize := 20
	res := make([]BossObj, nsize)
	for i := 0; i < nsize; i++ {
		b := BossObj{
			Platform: getPlatform(),
			Province: getProvince(),
			Isp:      getIsp(),
			Idc:      getIdc(),
			Ip:       randomdata.IpV4Address(),
			Hostname: randomdata.SillyName() + "_" + randomdata.StringNumberExt(4, "-", 2),
		}
		res[i] = b
	}
	slcB, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
	}
	f.Write(slcB)
}
