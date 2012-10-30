// Simple wrapper for rrdtool C library
package rrd

import (
	"strings"
	"time"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type Create struct {
	filename string
	start    time.Time
	step     time.Duration
	args     []string
}

func (c *Create) push(cmd, args []string) {
	c.args = append(c.args, strings.Join(append(cmd, args...), ":"))
}

// NewCreate returns new Create object. You need to call Save or Overwrite
// to really create database in filesystem.
func NewCreate(filename string, start time.Time, step time.Duration) *Create {

	return &Create{
		filename: filename,
		start:    start,
		step:     step,
	}
}

func (c *Create) DS(name, compute string, args ...string) {
	cmd := []string{"DS", name, compute}
	c.push(cmd, args)
}

func (c *Create) RRA(cf string, args ...string) {
	cmd := []string{"RRA", cf}
	c.push(cmd, args)
}

func (c *Create) Overwrite() error {

	return c.create()
}
