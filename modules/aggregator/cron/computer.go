package cron

import (
	"errors"
	"regexp"
	"strconv"
)

func compute(operands []string, operators []string, computeMode string, hostname string, valMap map[string]float64) (val float64, err error) {

	count := len(operands)
	if count == 0 {
		return val, errors.New("counter not found")
	}

	vals := queryOperands(operands, hostname, valMap)
	if len(vals) != count {
		return val, errors.New("value invalid")
	}

	sum := vals[0]
	for i, v := range vals[1:] {
		if operators[i] == "+" {
			sum += v
		} else {
			sum -= v
		}
	}

	if computeMode != "" {
		if compareSum(sum, computeMode) {
			val = 1
		}
	} else {
		val = sum
	}
	return val, nil
}

func compareSum(sum float64, computeMode string) bool {

	regMatch, _ := regexp.Compile(`([><=]+)([\d\.]+)`)
	match := regMatch.FindStringSubmatch(computeMode)

	mode := match[1]
	val, _ := strconv.ParseFloat(match[2], 64)

	switch {
	case mode == ">" && sum > val:
	case mode == "<" && sum < val:
	case mode == "=" && sum == val:
	case mode == ">=" && sum >= val:
	case mode == "<=" && sum <= val:
	default:
		return false
	}
	return true
}

func queryOperands(counters []string, endpoint string, valMap map[string]float64) []float64 {
	ret := []float64{}
	for _, counter := range counters {
		if v, ok := valMap[endpoint+counter]; ok {
			ret = append(ret, v)
		}
	}

	return ret
}
