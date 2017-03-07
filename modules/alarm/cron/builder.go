package cron

import (
	"fmt"

	"github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"[P%d][%s][%s][][%s %s %s %s %s%s%s][O%d %s]",
		event.Priority(),
		event.Status,
		event.Endpoint,
		event.Note(),
		event.Func(),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\nP%d\r\nEndpoint:%s\r\nMetric:%s\r\nTags:%s\r\n%s: %s%s%s\r\nNote:%s\r\nMax:%d, Current:%d\r\nTimestamp:%s\r\n%s\r\n",
		event.Status,
		event.Priority(),
		event.Endpoint,
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
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
