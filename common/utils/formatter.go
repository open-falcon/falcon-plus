package utils

import (
	"strconv"
	"strings"
)

func ReadableFloat(raw float64) string {
	val := strconv.FormatFloat(raw, 'f', 5, 64)
	if strings.Contains(val, ".") {
		val = strings.TrimRight(val, "0")
		val = strings.TrimRight(val, ".")
	}

	return val
}
