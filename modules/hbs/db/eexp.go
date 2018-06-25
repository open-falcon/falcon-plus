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

func QueryEExps() (ret []model.EExp, err error) {
	sql := "select id, exp, priority, max_step, note from eexp where pause=0"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("DB.Query failed", err)
		return ret, err
	}

	defer rows.Close()
	for rows.Next() {
		var ID int
		var exp string
		var priority int
		var max_step int
		var note string

		err = rows.Scan(
			&ID,
			&exp,
			&priority,
			&max_step,
			&note,
		)

		if err != nil {
			log.Println("parse result failed", err)
			continue
		}

		ee, err := parseEExp(exp)
		if err != nil {
			log.Println("parseEExp failed", err)
			continue
		}
		ee.ID = ID
		ee.Priority = priority
		ee.MaxStep = max_step
		ee.Note = note

		ret = append(ret, *ee)
	}

	return ret, nil
}

func parseFilter(s string) (*model.Filter, error) {
	filter := model.Filter{}
	var err error

	idxLeft := strings.Index(s, "(")
	idxRight := strings.Index(s, ")")
	if idxLeft == -1 || idxRight == -1 {
		err = errors.New("parse branket failed")
		return nil, err
	}

	filter.Func = strings.TrimSpace(s[:idxLeft])
	p := s[idxLeft+1 : idxRight]
	parts := strings.Split(p, ",")

	if filter.Func == "all" {
		if len(parts) != 2 {
			errmsg := fmt.Sprintf("func all parameter -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

		filter.Key = strings.TrimSpace(parts[0])
		buf := strings.TrimSpace(parts[1])

		splits := strings.Split(buf, "#")
		if len(splits) != 2 {
			errmsg := fmt.Sprintf("func all parameter -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

		filter.Limit, err = strconv.ParseUint(splits[1], 10, 64)
		if err != nil {
			errmsg := fmt.Sprintf("func all parameter -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

	} else if filter.Func == "count" {
		if len(parts) != 3 {
			errmsg := fmt.Sprintf("func count parameter -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

		filter.Key = strings.TrimSpace(parts[0])
		if filter.Key == "" {
			errmsg := fmt.Sprintf("func ago key -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

		filter.Ago, err = strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			errmsg := fmt.Sprintf("func ago parameter ago -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}

		filter.Hits, err = strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
		if err != nil {
			errmsg := fmt.Sprintf("func ago parameter hits -%s- is invalid", p)
			err = errors.New(errmsg)
			return nil, err
		}
	} else {
		err = errors.New(fmt.Sprintf("func -%s- not support", filter.Func))
		return nil, err
	}

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

	filter.Operator = buffer.String()

	remain = remain[strings.Index(remain, filter.Operator)+len(filter.Operator):]

	valueS := strings.TrimSpace(remain)
	if valueS == "" {
		err = errors.New("exp is invalid, value is empty")
		return nil, err
	}

	if valueS[0] == '"' {
		filter.RightValue = strings.Trim(valueS, `"`)
	} else {

		filter.RightValue, err = strconv.ParseFloat(valueS, 64)
		if err != nil {
			return nil, err
		}

	}

	return &filter, nil
}

func parseEExp(s string) (*model.EExp, error) {
	var err error
	ee := model.EExp{}
	ee.Filters = map[string]model.Filter{}

	s = strings.TrimSpace(s)
	if s == "" {
		err = errors.New("eexp is empty")
		return nil, err
	}

	keys := []string{}
	for _, filterS := range strings.Split(s, ";") {
		filterS = strings.TrimSpace(filterS)
		filter, err := parseFilter(filterS)
		if err != nil {
			log.Println("parseFilter failed", err)
		} else {
			keys = append(keys, filter.Key)
			ee.Filters[filter.Key] = *filter
		}
	}

	if len(keys) == 0 {
		err = errors.New("filters are invalid")
		return nil, err
	}

	sort.Sort(sort.StringSlice(keys))
	ee.Key = strings.Join(keys, ",")

	return &ee, nil
}
