package helper

import (
	"errors"
	"strconv"
)

func PageParser(page string, limit string) (p int, l int, err error) {
	p = -1
	l = -1
	if page != "" {
		p, err = strconv.Atoi(page)
		if err != nil {
			return
		}
		if limit != "" {
			l, err = strconv.Atoi(limit)
			if err != nil {
				return
			}
		} else {
			err = errors.New("You set page but skip limit params, please check your input")
			return
		}
		if p <= 0 || l <= 0 {
			err = errors.New("limit or page can not set to 0 or less than 0")
			return
		}
		if p == 1 {
			p = 0
		} else if p != 1 {
			p = ((p - 1) * l)
		}
	}
	return
}
