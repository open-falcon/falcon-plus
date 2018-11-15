// Copyright 2018 Xiaomi, Inc.
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

package event

import (
	"errors"
	"time"

	mevent "github.com/open-falcon/falcon-plus/modules/alarm-manager/model/event"
)

func checkEventInputs(inputs mevent.EventApiInputs) error {
	if inputs.UserName == "" && inputs.Uic == nil {
		if inputs.StartTime == 0 && inputs.EndTime == 0 {
			return errors.New("username or uic, start_time, end_time, You have to at least pick one on the request.")
		}
	}
	return nil
}

func TimeParamFilters(inputs mevent.EventApiInputs) mevent.EventApiInputs {
	if inputs.StartTime == 0 && inputs.EndTime == 0 {
		inputs.StartTime = time.Now().Unix() - 3600
		inputs.EndTime = time.Now().Unix()
	} else if inputs.StartTime == 0 {
		inputs.StartTime = time.Now().Unix() - 3600
	} else if inputs.EndTime == 0 {
		inputs.EndTime = time.Now().Unix()
	}
	return inputs
}

func LimitParamFilters(inputs mevent.EventApiInputs) mevent.EventApiInputs {
	if inputs.Limit == 0 || inputs.Limit >= 100 {
		inputs.Limit = 100
	}
	return inputs
}

func TimeQueryLimitFilters(inputs mevent.EventApiInputs) mevent.EventApiInputs {
	inputs = TimeParamFilters(inputs)
	return LimitParamFilters(inputs)
}
