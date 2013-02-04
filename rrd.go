// Simple wrapper for rrdtool C library
package rrd

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
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

// NewCreator returns new Creator object. You need to call Create to really
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
// If you specify args it saves them immediately.
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

type GraphInfo struct {
	Print         []string
	Width, Height uint
	Ymin, Ymax    float64
}

type Grapher struct {
	m               sync.Mutex
	title           string
	vlabel          string
	width, height   uint
	upperLimit      float64
	lowerLimit      float64
	rigid           bool
	altAutoscale    bool
	altAutoscaleMin bool
	altAutoscaleMax bool
	noGridFit       bool

	logarithmic bool

	noLegend bool

	lazy bool

	color string

	slopeMode bool

	watermark   string
	base        uint
	imageFormat string
	interlaced  bool

	args []string
}

func NewGrapher() *Grapher {
	return &Grapher{
		upperLimit: -math.MaxFloat64,
		lowerLimit: math.MaxFloat64,
	}
}

func (g *Grapher) SetTitle(title string) {
	g.title = title
}

func (g *Grapher) SetVLabel(vlabel string) {
	g.vlabel = vlabel
}

func (g *Grapher) SetSize(width, height uint) {
	g.width = width
	g.height = height
}

func (g *Grapher) SetLowerLimit(limit float64) {
	g.lowerLimit = limit
}

func (g *Grapher) SetUpperLimit(limit float64) {
	g.upperLimit = limit
}

func (g *Grapher) SetRigid() {
	g.rigid = true
}

func (g *Grapher) SetAltAutoscale() {
	g.altAutoscale = true
}

func (g *Grapher) SetAltAutoscaleMin() {
	g.altAutoscaleMin = true
}

func (g *Grapher) SetAltAutoscaleMax() {

	g.altAutoscaleMax = true
}

func (g *Grapher) SetNoGridFit() {
	g.noGridFit = true
}

func (g *Grapher) SetLogarithmic() {
	g.logarithmic = true
}

func (g *Grapher) SetNoLegend() {
	g.noLegend = true
}

func (g *Grapher) SetLazy() {
	g.lazy = true
}

func (g *Grapher) SetColor(colortag, color string) {
	g.color = colortag + "#" + color
}

func (g *Grapher) SetSlopeMode() {
	g.slopeMode = true
}

func (g *Grapher) SetImageFormat(format string) {
	g.imageFormat = format
}

func (g *Grapher) SetInterlaced() {
	g.interlaced = true
}

func (g *Grapher) SetBase(base uint) {
	g.base = base
}

func (g *Grapher) SetWatermark(watermark string) {
	g.watermark = watermark
}

func (g *Grapher) push(cmd string, options []string) {
	if len(options) > 0 {
		cmd += ":" + strings.Join(options, ":")
	}
	g.args = append(g.args, cmd)
}

func (g *Grapher) Def(vname, rrdfile, dsname, cf string, options ...string) {
	g.push(
		fmt.Sprintf("DEF:%s=%s:%s:%s", vname, rrdfile, dsname, cf),
		options,
	)
}

func (g *Grapher) VDef(vname, rpn string) {
	g.push("VDEF:"+vname+"="+rpn, nil)
}

func (g *Grapher) CDef(vname, rpn string) {
	g.push("CDEF:"+vname+"="+rpn, nil)
}

func (g *Grapher) Print(vname, format string) {
	g.push("PRINT:"+vname+":"+format, nil)
}

func (g *Grapher) PrintT(vname, format string) {
	g.push("PRINT:"+vname+":"+format+":strftime", nil)
}
func (g *Grapher) GPrint(vname, format string) {
	g.push("GPRINT:"+vname+":"+format, nil)
}

func (g *Grapher) GPrintT(vname, format string) {
	g.push("GPRINT:"+vname+":"+format+":strftime", nil)
}

func (g *Grapher) Comment(s string) {
	g.push("COMMENT:"+s, nil)
}

func (g *Grapher) VRule(t interface{}, color string, options ...string) {
	if v, ok := t.(time.Time); ok {
		t = v.Unix()
	}
	vr := fmt.Sprintf("VRULE:%s#%s", t, color)
	g.push(vr, options)
}

func (g *Grapher) HRule(value, color string, options ...string) {
	hr := "HRULE:" + value + "#" + color
	g.push(hr, options)
}

func (g *Grapher) Line(width float32, value, color string, options ...string) {
	line := fmt.Sprintf("LINE%f:%s", width, value)
	if color != "" {
		line += "#" + color
	}
	g.push(line, options)
}

func (g *Grapher) Area(value, color string, options ...string) {
	area := "AREA:" + value
	if color != "" {
		area += "#" + color
	}
	g.push(area, options)
}

func (g *Grapher) Tick(vname, color string, options ...string) {
	tick := "TICK:" + vname
	if color != "" {
		tick += "#" + color
	}
	g.push(tick, options)
}

func (g *Grapher) Shift(vname string, offset interface{}) {
	if v, ok := offset.(time.Duration); ok {
		offset = int64((v + time.Second/2) / time.Second)
	}
	shift := fmt.Sprintf("SHIFT:%s:%s", offset)
	g.push(shift, nil)
}

func (g *Grapher) TextAlign(align string) {
	g.push("TEXTALIGN:"+align, nil)
}

// Graph returns GraphInfo and image as []byte or error
func (g *Grapher) Graph(start, end time.Time) (GraphInfo, []byte, error) {
	return g.graph("-", start, end)
}

// SaveGraph saves image to file and returns GraphInfo or error
func (g *Grapher) SaveGraph(filename string, start, end time.Time) (GraphInfo, error) {
	gi, _, err := g.graph(filename, start, end)
	return gi, err
}
