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
			result = fmt.Sprintf("'%v'", a)
		} else {
			result = fmt.Sprintf("%v,'%v'", result, a)
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
