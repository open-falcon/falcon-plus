package db

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/open-falcon/falcon-plus/common/model"
)

func QueryEExpressions() (ret []*model.EExpression, err error) {
	sql := "select id, filters, conditions, priority, max_step, note from eexp where pause=0"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("DB.Query failed", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var ID int
		var filters string
		var conditions string
		var priority int
		var max_step int
		var note string

		err = rows.Scan(
			&ID,
			&filters,
			&conditions,
			&priority,
			&max_step,
			&note,
		)

		if err != nil {
			log.Println("parse result failed", err)
			continue
		}

		ee, err := parseEExpression(filters, conditions)
		if err != nil {
			log.Println("parseEExpression failed", err)
			continue
		}
		ee.ID = ID
		ee.Priority = priority
		ee.MaxStep = max_step
		ee.Note = note

		ret = append(ret, ee)
	}

	return ret, nil
}

func parseCond(s string) (*model.Condition, error) {
	c := model.Condition{}
	var err error

	idxLeft := strings.Index(s, "(")
	idxRight := strings.Index(s, ")")
	if idxLeft == -1 || idxRight == -1 {
		err = errors.New("parse branket failed")
		return nil, err
	}

	c.Func = s[:idxLeft]
	p := s[idxLeft+1 : idxRight]

	parts := strings.Split(p, ",")
	if len(parts) != 2 {
		err = errors.New("parse parameter failed")
		return nil, err
	}
	c.Metric = strings.TrimSpace(parts[0])
	c.Parameter = strings.TrimSpace(parts[1])

	parts = strings.Split(c.Parameter, "#")
	if len(parts) != 2 {
		err = errors.New(fmt.Sprintf("parameter -%s- is invalid", c.Parameter))
		log.Println("parse parameter failed", err)
		return nil, err
	}
	limit, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		err = errors.New(fmt.Sprintf("parameter -%s- is invalid", c.Parameter))
		log.Println("parse parameter limit failed", err)
		return nil, err
	}
	c.Limit = int(limit)

	remain := strings.TrimSpace(s[idxRight+1:])

	var buffer bytes.Buffer
	var chr rune
	for _, chr = range remain {
		switch chr {
		case '=':
			{
				buffer.WriteRune(chr)
			}
		case '>':
			{

				buffer.WriteRune(chr)
			}

		case '<':
			{

				buffer.WriteRune(chr)
			}
		case ' ':
			{

			}
		default:
			{
				break
			}
		}
	}

	c.Operator = buffer.String()

	remain = remain[strings.Index(remain, c.Operator)+len(c.Operator):]

	valueI := strings.TrimSpace(remain)
	value, err := strconv.ParseFloat(valueI, 64)
	if err != nil {
		return nil, err
	}

	c.RightValue = value

	return &c, nil
}

func parseEExpression(filter string, conds string) (*model.EExpression, error) {
	var err error
	ee := model.EExpression{}
	ee.Filters = map[string]string{}

	filter = strings.TrimSpace(filter)
	if filter == "" {
		err = errors.New("filter is empty")
		return nil, err
	}

	idxLeft := strings.Index(filter, "(")
	idxRight := strings.Index(filter, ")")
	if idxLeft == -1 || idxRight == -1 {
		err = errors.New("filter is empty")
		return nil, err
	}

	ee.Func = filter[:idxLeft]

	kvPairs := strings.Split(filter[idxLeft+1:idxRight], " ")
	for _, kvPair := range kvPairs {
		kvSlice := strings.Split(kvPair, "=")
		if len(kvSlice) != 2 {
			continue
		}

		key := kvSlice[0]
		ee.Filters[key] = kvSlice[1]
	}

	metrics := []string{}
	for _, cond := range strings.Split(conds, ";") {
		cond = strings.TrimSpace(cond)
		c, err := parseCond(cond)
		if err != nil {
			log.Println("parse condition failed", err)
		} else {
			metrics = append(metrics, c.Metric)
			ee.Conditions = append(ee.Conditions, *c)
		}
	}

	if len(metrics) == 0 {
		err = errors.New("conditions are invalid")
		return nil, err
	}

	sort.Sort(sort.StringSlice(metrics))
	ee.Metric = strings.Join(metrics, ",")

	return &ee, nil
}
