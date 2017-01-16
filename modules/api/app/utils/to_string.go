package utils

import (
	"errors"
	"fmt"
)

func ArrIntToString(arr []int) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}

func ArrInt64ToString(arr []int64) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}
