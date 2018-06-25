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

package cron

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func BuildCommonSMSContent(event *model.Event) string {
	var leftValue string
	var rightValue string
	pushedTags := map[string]string{}

	for k, v := range event.PushedTags {
		pushedTags[k] = fmt.Sprintf("%v", v)
	}

	switch event.LeftValue.(type) {
	case float64:
		{
			leftValue = utils.ReadableFloat(event.LeftValue.(float64))
			rightValue = utils.ReadableFloat(event.RightValue())
		}
	case string:
		{
			leftValue = event.LeftValue.(string)
			rightValue = fmt.Sprintf("%v", event.RightValue())

		}

	}
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %v%s%v][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		pushedTags,
		leftValue,
		event.Operator(),
		rightValue,
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonIMContent(event *model.Event) string {
	var leftValue string
	var rightValue string
	pushedTags := map[string]string{}

	for k, v := range event.PushedTags {
		pushedTags[k] = fmt.Sprintf("%v", v)
	}

	switch event.LeftValue.(type) {
	case float64:
		{
			leftValue = utils.ReadableFloat(event.LeftValue.(float64))
			rightValue = utils.ReadableFloat(event.RightValue())
		}
	case string:
		{
			leftValue = event.LeftValue.(string)
			rightValue = fmt.Sprintf("%v", event.RightValue())

		}

	}

	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		pushedTags,
		leftValue,
		event.Operator(),
		rightValue,
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	var leftValue string
	var rightValue string
	pushedTags := map[string]string{}

	for k, v := range event.PushedTags {
		pushedTags[k] = fmt.Sprintf("%v", v)
	}

	switch event.LeftValue.(type) {
	case float64:
		{
			leftValue = utils.ReadableFloat(event.LeftValue.(float64))
			rightValue = utils.ReadableFloat(event.RightValue())
		}
	case string:
		{
			leftValue = event.LeftValue.(string)
			rightValue = fmt.Sprintf("%v", event.RightValue())

		}

	}

	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\nP%d\r\nEndpoint:%s\r\nMetric:%s\r\nTags:%s\r\n%s: %s%s%s\r\nNote:%s\r\nMax:%d, Current:%d\r\nTimestamp:%s\r\n%s\r\n",
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		pushedTags,
		event.Func(),
		leftValue,
		event.Operator(),
		rightValue,
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateIMContent(event *model.Event) string {
	return BuildCommonIMContent(event)
}
