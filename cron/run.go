package cron

import (
	"fmt"
	"github.com/open-falcon/aggregator/g"
	"github.com/open-falcon/sdk/portal"
	"github.com/open-falcon/sdk/sender"
	"log"
	"strconv"
	"strings"
	"time"
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

	numeratorOperands, numeratorOperators := parse(numeratorStr, needComputeNumerator)
	denominatorOperands, denominatorOperators := parse(denominatorStr, needComputeDenominator)

	if !operatorsValid(numeratorOperators) || !operatorsValid(denominatorOperators) {
		log.Println("[W] operators invalid", item)
		return
	}

	hostnames, err := portal.Hostnames(fmt.Sprintf("%d", item.GroupId))
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
		var (
			numeratorValid   = true
			denominatorValid = true
		)

		if needComputeNumerator {
			numeratorVal, numeratorValid = compute(numeratorOperands, numeratorOperators, hostname, valueMap)
			if !numeratorValid && debug {
				log.Printf("[W] [hostname:%s] [numerator:%s] invalid or not found", hostname, item.Numerator)
			}
		}

		if needComputeDenominator {
			denominatorVal, denominatorValid = compute(denominatorOperands, denominatorOperators, hostname, valueMap)
			if !denominatorValid && debug {
				log.Printf("[W] [hostname:%s] [denominator:%s] invalid or not found", hostname, item.Denominator)
			}
		}

		if numeratorValid && denominatorValid {
			numerator += numeratorVal
			denominator += denominatorVal
			validCount += 1
		}
	}

	if !needComputeNumerator {
		if numeratorStr == "$#" {
			numerator = float64(validCount)
		} else {
			numerator, err = strconv.ParseFloat(numeratorStr, 64)
			if err != nil {
				log.Printf("[E] strconv.ParseFloat(%s) fail %v", numeratorStr, item)
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
				log.Printf("[E] strconv.ParseFloat(%s) fail %v", denominatorStr, item)
				return
			}
		}
	}

	if denominator == 0 {
		log.Println("[W] denominator == 0", item)
		return
	}

	sender.Push(item.Endpoint, item.Metric, item.Tags, numerator/denominator, item.DsType, int64(item.Step))
}

func parse(expression string, needCompute bool) (operands []string, operators []uint8) {
	if !needCompute {
		return
	}

	// e.g. $(cpu.busy)+$(cpu.idle)-$(cpu.nice)-$(cpu.guest)
	//      xx          --          --          --         x
	newExpression := expression[2 : len(expression)-1]
	arr := strings.Split(newExpression, "$(")
	count := len(arr)
	if count == 1 {
		operands = append(operands, arr[0])
		return
	}

	if count > 1 {
		for i := 0; i < count; i++ {
			item := arr[i]
			length := len(item)
			if i == count-1 {
				operands = append(operands, item)
				continue
			}
			operators = append(operators, item[length-1])
			operands = append(operands, item[0:length-2])
		}
	}

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
	if strings.Contains(val, "$(") {
		return true
	}

	return false
}

func expressionValid(val string) bool {
	// use chinese character?
	if strings.Contains(val, "（") || strings.Contains(val, "）") {
		return false
	}

	return true
}

func operatorsValid(ops []uint8) bool {
	count := len(ops)
	for i := 0; i < count; i++ {
		if ops[i] != '+' && ops[i] != '-' {
			return false
		}
	}
	return true
}
