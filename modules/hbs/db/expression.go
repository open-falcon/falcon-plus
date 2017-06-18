package db

import (
	"fmt"
	"github.com/open-falcon/falcon-plus/common/model"
	"log"
	"strings"
)

func QueryExpressions() (ret []*model.Expression, err error) {
	sql := "select id, expression, func, op, right_value, max_step, priority, note, action_id from expression where action_id>0 and pause=0"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		e := model.Expression{}
		var exp string
		err = rows.Scan(
			&e.Id,
			&exp,
			&e.Func,
			&e.Operator,
			&e.RightValue,
			&e.MaxStep,
			&e.Priority,
			&e.Note,
			&e.ActionId,
		)

		if err != nil {
			log.Println("WARN:", err)
			continue
		}

		e.Metric, e.Tags, err = parseExpression(exp)
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}

		ret = append(ret, &e)
	}

	return ret, nil
}

func parseExpression(exp string) (metric string, tags map[string]string, err error) {
	left := strings.Index(exp, "(")
	right := strings.Index(exp, ")")
	tagStrs := strings.TrimSpace(exp[left+1 : right])

	arr := strings.Fields(tagStrs)
	if len(arr) < 2 {
		err = fmt.Errorf("tag not enough. exp: %s", exp)
		return
	}

	tags = make(map[string]string)
	for _, item := range arr {
		kv := strings.Split(item, "=")
		if len(kv) != 2 {
			err = fmt.Errorf("parse %s fail", exp)
			return
		}
		tags[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	metric, exists := tags["metric"]
	if !exists {
		err = fmt.Errorf("no metric give of %s", exp)
		return
	}

	delete(tags, "metric")
	return
}
