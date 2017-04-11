package event

import (
	"fmt"
	"time"

	"database/sql"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
	coommonModel "github.com/open-falcon/falcon-plus/common/model"
	"github.com/open-falcon/falcon-plus/common/utils"
)

const timeLayout = "2006-01-02 15:04:05"

func insertEvent(q orm.Ormer, eve *coommonModel.Event) (res interface{}, err error) {
	var status int
	if status = 0; eve.Status == "OK" {
		status = 1
	}
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		status,
		timestamp
	) VALUES(?,?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		eve.Id,
		eve.CurrentStep,
		fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		status,
		time.Unix(eve.EventTime, 0).Format(timeLayout),
	).Exec()
	return
}

func InsertEvent(eve *coommonModel.Event) {
	q := orm.NewOrm()
	var event []EventCases
	q.Raw("select * from event_cases where id = ?", eve.Id).QueryRows(&event)
	var sqlLog sql.Result
	var errRes error
	log.Debugf("events: %v", eve)
	log.Debugf("expression is null: %v", eve.Expression == nil)
	if len(event) == 0 {
		//create cases
		sqltemplete := `INSERT INTO event_cases (
					id,
					endpoint,
					metric,
					func,
					cond,
					note,
					max_step,
					current_step,
					priority,
					status,
					timestamp,
					update_at,
					tpl_creator,
					expression_id,
					strategy_id,
					template_id
					) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		sqlLog, errRes = q.Raw(
			sqltemplete,
			eve.Id,
			eve.Endpoint,
			counterGen(eve.Metric(), utils.SortedTags(eve.PushedTags)),
			eve.Func(),
			//cond
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Strategy.Note,
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Priority(),
			eve.Status,
			//start_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			//update_at
			time.Unix(eve.EventTime, 0).Format(timeLayout),
			eve.Strategy.Tpl.Creator,
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId()).Exec()

	} else {
		sqltemplete := `UPDATE event_cases SET
				update_at = ?,
				max_step = ?,
				current_step = ?,
				note = ?,
				cond = ?,
				status = ?,
				func = ?,
				priority = ?,
				tpl_creator = ?,
				expression_id = ?,
				strategy_id = ?,
				template_id = ?`
		//reopen case
		if event[0].ProcessStatus == "resolved" || event[0].ProcessStatus == "ignored" {
			sqltemplete = fmt.Sprintf("%v ,process_status = '%s', process_note = %d", sqltemplete, "unresolved", 0)
		}

		if eve.CurrentStep == 1 {
			//update start time of cases
			sqltemplete = fmt.Sprintf("%v ,timestamp = ? WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Strategy.Note,
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				eve.Strategy.Tpl.Creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.Id,
			).Exec()
		} else {
			sqltemplete = fmt.Sprintf("%v WHERE id = ?", sqltemplete)
			sqlLog, errRes = q.Raw(
				sqltemplete,
				time.Unix(eve.EventTime, 0).Format(timeLayout),
				eve.MaxStep(),
				eve.CurrentStep,
				eve.Strategy.Note,
				fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
				eve.Status,
				eve.Func(),
				eve.Priority(),
				eve.Strategy.Tpl.Creator,
				eve.ExpressionId(),
				eve.StrategyId(),
				eve.TplId(),
				eve.Id,
			).Exec()
		}
	}
	log.Debug(fmt.Sprintf("%v, %v", sqlLog, errRes))
	//insert case
	res2, err := insertEvent(q, eve)
	log.Debugf("insert event to db response:%v, error:%v", res2, err)
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}
