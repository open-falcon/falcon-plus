package cron

func compute(operands []string, operators []uint8, hostname string, valMap map[string]float64) (val float64, valid bool) {
	count := len(operands)
	if count == 0 {
		return val, false
	}

	vals := queryOperands(operands, hostname, valMap)
	if len(vals) != count {
		return val, false
	}

	val = vals[0]

	for i := 1; i < count; i++ {
		if operators[i] == '+' {
			val += vals[i]
		} else {
			val -= vals[i]
		}
	}

	return val, true
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
