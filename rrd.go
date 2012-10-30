// Simple wrapper for rrdtool C library
package rrd

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type cstring []byte

func newCstring(s string) cstring {
	cs := make(cstring, len(s)+1)
	copy(cs, s)
	return cs
}

func (cs cstring) p() unsafe.Pointer {
	if len(cs) == 0 {
		return nil
	}
	return unsafe.Pointer(&cs[0])
}

func (cs cstring) String() string {
	return string(cs[:len(cs)-1])
}

func join(args []interface{}) string {
	sa := make([]string, len(args))
	for i, a := range args {
		var s string
		switch v := a.(type) {
		case time.Time:
			s = fmt.Sprint(v.Unix())
		default:
			s = fmt.Sprint(v)
		}
		sa[i] = s
	}
	return strings.Join(sa, ":")
}

type Creator struct {
	filename string
	start    time.Time
	step     uint
	args     []string
}

// NewCreate returns new Creator object. You need to call Create to really
// create database file.
//	filename - name of database file
//	start    - don't accept any data timed before or at time specified
//	step     - base interval in seconds with which data will be fed into RRD
func NewCreator(filename string, start time.Time, step uint) *Creator {
	return &Creator{
		filename: filename,
		start:    start,
		step:     step,
	}
}

func (c *Creator) DS(name, compute string, args ...interface{}) {
	c.args = append(c.args, "DS:"+name+":"+compute+":"+join(args))
}

func (c *Creator) RRA(cf string, args ...interface{}) {
	c.args = append(c.args, "RRA:"+cf+":"+join(args))
}

func (c *Creator) Create(overwrite bool) error {
	if !overwrite {
		f, err := os.OpenFile(
			c.filename,
			os.O_WRONLY|os.O_CREATE|os.O_EXCL,
			0666,
		)
		if err != nil {
			return err
		}
		f.Close()
	}
	return c.create()
}

// Use cstring and unsafe.Pointer to avoid alocations for C calls

type Updater struct {
	filename cstring
	template cstring

	args []unsafe.Pointer
}

func NewUpdater(filename string) *Updater {
	return &Updater{filename: newCstring(filename)}
}

func (u *Updater) SetTemplate(dsName ...string) {
	u.template = newCstring(strings.Join(dsName, ":"))
}

// Cache chaches data for later save using Update(). Use it to avoid
// open/read/write/close for every update.
func (u *Updater) Cache(args ...interface{}) {
	u.args = append(u.args, newCstring(join(args)).p())
}

// Update saves data in RRDB.
// Without args Update saves all subsequent updates buffered by Cache method.
// If you specify args it saves them immediately (this is thread-safe
// operation).
func (u *Updater) Update(args ...interface{}) error {
	if len(args) != 0 {
		a := make([]unsafe.Pointer, 1)
		a[0] = newCstring(join(args)).p()
		return u.update(a)
	} else if len(u.args) != 0 {
		err := u.update(u.args)
		u.args = nil
		return err
	}
	return nil
}
