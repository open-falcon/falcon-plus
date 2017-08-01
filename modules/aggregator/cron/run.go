package cron

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/open-falcon/falcon-plus/common/sdk/sender"
	"github.com/open-falcon/falcon-plus/modules/aggregator/g"
	"github.com/open-falcon/falcon-plus/modules/aggregator/sdk"
)

func WorkerRun(item *g.Cluster) {
	debug := g.Config().Debug

	numeratorStr := cleanParam(item.Numerator)
	denominatorStr := cleanParam(item.Denominator)

	if !expressionValid(numeratorStr) || !expressionValid(denominatorStr) {
		log.Println("[W] invalid numerator or denominator", item)
		return
	}

	needComputeNumerator := needCompute(numeratorStr)
	needComputeDenominator := needCompute(denominatorStr)

	if !needComputeNumerator && !needComputeDenominator {
		log.Println("[W] no need compute", item)
		return
	}

	numeratorOperands, numeratorOperators, numeratorComputeMode := parse(numeratorStr, needComputeNumerator)
	denominatorOperands, denominatorOperators, denominatorComputeMode := parse(denominatorStr, needComputeDenominator)

	if !operatorsValid(numeratorOperators) || !operatorsValid(denominatorOperators) {
		log.Println("[W] operators invalid", item)
		return
	}

	hostnames, err := sdk.HostnamesByID(item.GroupId)
	if err != nil || len(hostnames) == 0 {
		return
	}

	now := time.Now().Unix()

	valueMap, err := queryCounterLast(numeratorOperands, denominatorOperands, hostnames, now-int64(item.Step*2), now)
	if err != nil {
		log.Println("[E]", err, item)
		return
	}

	var numerator, denominator float64
	var validCount int

	for _, hostname := range hostnames {
		var numeratorVal, denominatorVal float64
		var err error

		if needComputeNumerator {
			numeratorVal, err = compute(numeratorOperands, numeratorOperators, numeratorComputeMode, hostname, valueMap)

			if debug && err != nil {
				log.Printf("[W] [hostname:%s] [numerator:%s] id:%d, err:%v", hostname, item.Numerator, item.Id, err)
			} else if debug {
				log.Printf("[D] [hostname:%s] [numerator:%s] id:%d, value:%0.4f", hostname, item.Numerator, item.Id, numeratorVal)
			}

			if err != nil {
				continue
			}
		}

		if needComputeDenominator {
			denominatorVal, err = compute(denominatorOperands, denominatorOperators, denominatorComputeMode, hostname, valueMap)

			if debug && err != nil {
				log.Printf("[W] [hostname:%s] [denominator:%s] id:%d, err:%v", hostname, item.Denominator, item.Id, err)
			} else if debug {
				log.Printf("[D] [hostname:%s] [denominator:%s] id:%d, value:%0.4f", hostname, item.Denominator, item.Id, denominatorVal)
			}

			if err != nil {
				continue
			}
		}

		if debug {
			log.Printf("[D] hostname:%s  numerator:%0.4f  denominator:%0.4f  per:%0.4f\n", hostname, numeratorVal, denominatorVal, numeratorVal/denominatorVal)
		}
		numerator += numeratorVal
		denominator += denominatorVal
		validCount += 1
	}

	if !needComputeNumerator {
		if numeratorStr == "$#" {
			numerator = float64(validCount)
		} else {
			numerator, err = strconv.ParseFloat(numeratorStr, 64)
			if err != nil {
				log.Printf("[E] strconv.ParseFloat(%s) fail %v, id:%d", numeratorStr, item.Id)
				return
			}
		}
	}

	if !needComputeDenominator {
		if denominatorStr == "$#" {
			denominator = float64(validCount)
		} else {
			denominator, err = strconv.ParseFloat(denominatorStr, 64)
			if err != nil {
				log.Printf("[E] strconv.ParseFloat(%s) fail %v, id:%d", denominatorStr, item.Id)
				return
			}
		}
	}

	if denominator == 0 {
		log.Println("[W] denominator == 0, id:", item.Id)
		return
	}

	if validCount == 0 {
		log.Println("[W] validCount == 0, id:", item.Id)
		return
	}

	if debug {
		log.Printf("[D] hostname:all  numerator:%0.4f  denominator:%0.4f  per:%0.4f\n", numerator, denominator, numerator/denominator)
	}
	sender.Push(item.Endpoint, item.Metric, item.Tags, numerator/denominator, item.DsType, int64(item.Step))
}

func parse(expression string, needCompute bool) (operands []string, operators []string, computeMode string) {
	if !needCompute {
		return
	}

	// e.g. $(cpu.busy)
	// e.g. $(cpu.busy)+$(cpu.idle)-$(cpu.nice)
	// e.g. $(cpu.busy)>=80
	// e.g. ($(cpu.busy)+$(cpu.idle)-$(cpu.nice))>80

	splitCounter, _ := regexp.Compile(`[\$\(\)]+`)
	items := splitCounter.Split(expression, -1)

	count := len(items)
	for i, val := range items[1 : count-1] {
		if i%2 == 0 {
			operands = append(operands, val)
		} else {
			operators = append(operators, val)
		}
	}
	computeMode = items[count-1]

	return
}

func cleanParam(val string) string {
	val = strings.TrimSpace(val)
	val = strings.Replace(val, " ", "", -1)
	val = strings.Replace(val, "\r", "", -1)
	val = strings.Replace(val, "\n", "", -1)
	val = strings.Replace(val, "\t", "", -1)
	return val
}

// $#
// 200
// $(cpu.busy) + $(cpu.idle)
func needCompute(val string) bool {
	return strings.Contains(val, "$(")
}

func expressionValid(val string) bool {
	// use chinese character?

	if strings.Contains(val, "（") || strings.Contains(val, "）") {
		return false
	}

	if val == "$#" {
		return true
	}

	// e.g. $(cpu.busy)
	// e.g. $(cpu.busy)+$(cpu.idle)-$(cpu.nice)
	matchMode0 := `^(\$\([^\(\)]+\)[+-])*\$\([^\(\)]+\)$`
	if ok, err := regexp.MatchString(matchMode0, val); err == nil && ok {
		return true
	}

	// e.g. $(cpu.busy)>=80
	matchMode1 := `^\$\([^\(\)]+\)(>|=|<|>=|<=)\d+(\.\d+)?$`
	if ok, err := regexp.MatchString(matchMode1, val); err == nil && ok {
		return true
	}

	// e.g. ($(cpu.busy)+$(cpu.idle)-$(cpu.nice))>80
	matchMode2 := `^\((\$\([^\(\)]+\)[+-])*\$\([^\(\)]+\)\)(>|=|<|>=|<=)\d+(\.\d+)?$`
	if ok, err := regexp.MatchString(matchMode2, val); err == nil && ok {
		return true
	}

	// e.g. 纯数字
	matchMode3 := `^\d+$`
	if ok, err := regexp.MatchString(matchMode3, val); err == nil && ok {
		return true
	}

	return false
}

func operatorsValid(ops []string) bool {
	count := len(ops)
	for i := 0; i < count; i++ {
		if ops[i] != "+" && ops[i] != "-" {
			return false
		}
	}
	return true
}
