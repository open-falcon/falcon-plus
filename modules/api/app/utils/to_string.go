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

func ArrIntToStringMust(arr []int) (result string) {
	result, _ = ArrIntToString(arr)
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

func ArrInt64ToStringMust(arr []int64) (result string) {
	result, _ = ArrInt64ToString(arr)
	return
}

func ArrStringsToString(arr []string) (result string, err error) {
	result = ""
	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("\"%v\"", a)
		} else {
			result = fmt.Sprintf("%v,\"%v\"", result, a)
		}
	}
	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}
	return
}

func ArrStringsToStringMust(arr []string) (result string) {
	result, _ = ArrStringsToString(arr)
	return
}
