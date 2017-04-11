package g

import (
	"fmt"
	cmodel "github.com/open-falcon/falcon-plus/common/model"
)

func Link(event *cmodel.Event) string {
	tplId := event.TplId()
	if tplId != 0 {
		return fmt.Sprintf("%s/portal/template/view/%d", Config().Api.Dashboard, tplId)
	}

	eid := event.ExpressionId()
	if eid != 0 {
		return fmt.Sprintf("%s/portal/expression/view/%d", Config().Api.Dashboard, eid)
	}

	return ""
}
