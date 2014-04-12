package rrd

/*
#include <stdlib.h>
#include <rrd.h>
#include "rrdfunc.h"
#cgo LDFLAGS: -lrrd_th
*/
import "C"
import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var mutex sync.Mutex

func makeArgs(args []string) []*C.char {
	ret := make([]*C.char, len(args))
	for i, s := range args {
		ret[i] = C.CString(s)
	}
	return ret
}

func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func freeArgs(cArgs []*C.char) {
	for _, s := range cArgs {
		freeCString(s)
	}
}

func makeError(e *C.char) error {
	var null *C.char
	if e == null {
		return nil
	}
	defer freeCString(e)
	return Error(C.GoString(e))
}

func (c *Creator) create() error {
	filename := C.CString(c.filename)
	defer freeCString(filename)
	args := makeArgs(c.args)
	defer freeArgs(args)

	e := C.rrdCreate(
		filename,
		C.ulong(c.step),
		C.time_t(c.start.Unix()),
		C.int(len(args)),
		&args[0],
	)
	return makeError(e)
}

func (u *Updater) update(args []unsafe.Pointer) error {
	e := C.rrdUpdate(
		(*C.char)(u.filename.p()),
		(*C.char)(u.template.p()),
		C.int(len(args)),
		(**C.char)(unsafe.Pointer(&args[0])),
	)
	return makeError(e)
}

var (
	graphv = C.CString("graphv")
	xport  = C.CString("xport")

	oStart           = C.CString("-s")
	oEnd             = C.CString("-e")
	oTitle           = C.CString("-t")
	oVlabel          = C.CString("-v")
	oWidth           = C.CString("-w")
	oHeight          = C.CString("-h")
	oUpperLimit      = C.CString("-u")
	oLowerLimit      = C.CString("-l")
	oRigid           = C.CString("-r")
	oAltAutoscale    = C.CString("-A")
	oAltAutoscaleMin = C.CString("-J")
	oAltAutoscaleMax = C.CString("-M")
	oNoGridFit       = C.CString("-N")

	oLogarithmic   = C.CString("-o")
	oUnitsExponent = C.CString("-X")
	oUnitsLength   = C.CString("-L")

	oRightAxis      = C.CString("--right-axis")
	oRightAxisLabel = C.CString("--right-axis-label")

	oDaemon = C.CString("--daemon")

	oNoLegend = C.CString("-g")

	oLazy = C.CString("-z")

	oColor = C.CString("-c")

	oSlopeMode   = C.CString("-E")
	oImageFormat = C.CString("-a")
	oInterlaced  = C.CString("-i")

	oBase      = C.CString("-b")
	oWatermark = C.CString("-W")

	oStep    = C.CString("--step")
	oMaxRows = C.CString("-m")
)

func ftoa(f float64) string {
	return strconv.FormatFloat(f, 'e', 10, 64)
}

func ftoc(f float64) *C.char {
	return C.CString(ftoa(f))
}

func i64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func i64toc(i int64) *C.char {
	return C.CString(i64toa(i))
}

func u64toa(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func u64toc(u uint64) *C.char {
	return C.CString(u64toa(u))
}
func itoa(i int) string {
	return i64toa(int64(i))
}

func itoc(i int) *C.char {
	return i64toc(int64(i))
}

func utoa(u uint) string {
	return u64toa(uint64(u))
}

func utoc(u uint) *C.char {
	return u64toc(uint64(u))
}

func (g *Grapher) makeArgs(filename string, start, end time.Time) []*C.char {
	args := []*C.char{
		graphv, C.CString(filename),
		oStart, i64toc(start.Unix()),
		oEnd, i64toc(end.Unix()),
		oTitle, C.CString(g.title),
		oVlabel, C.CString(g.vlabel),
	}
	if g.width != 0 {
		args = append(args, oWidth, utoc(g.width))
	}
	if g.height != 0 {
		args = append(args, oHeight, utoc(g.height))
	}
	if g.upperLimit != -math.MaxFloat64 {
		args = append(args, oUpperLimit, ftoc(g.upperLimit))
	}
	if g.lowerLimit != math.MaxFloat64 {
		args = append(args, oLowerLimit, ftoc(g.lowerLimit))
	}
	if g.rigid {
		args = append(args, oRigid)
	}
	if g.altAutoscale {
		args = append(args, oAltAutoscale)
	}
	if g.altAutoscaleMax {
		args = append(args, oAltAutoscaleMax)
	}
	if g.altAutoscaleMin {
		args = append(args, oAltAutoscaleMin)
	}
	if g.noGridFit {
		args = append(args, oNoGridFit)
	}
	if g.logarithmic {
		args = append(args, oLogarithmic)
	}
	if g.unitsExponent != minInt {
		args = append(
			args,
			oUnitsExponent, itoc(g.unitsExponent),
		)
	}
	if g.unitsLength != 0 {
		args = append(
			args,
			oUnitsLength, utoc(g.unitsLength),
		)
	}
	if g.rightAxisScale != 0 {
		args = append(
			args,
			oRightAxis,
			C.CString(ftoa(g.rightAxisScale)+":"+ftoa(g.rightAxisShift)),
		)
	}
	if g.rightAxisLabel != "" {
		args = append(
			args,
			oRightAxisLabel, C.CString(g.rightAxisLabel),
		)
	}
	if g.noLegend {
		args = append(args, oNoLegend)
	}
	if g.lazy {
		args = append(args, oLazy)
	}
	if g.color != "" {
		args = append(args, oColor, C.CString(g.color))
	}
	if g.slopeMode {
		args = append(args, oSlopeMode)
	}
	if g.imageFormat != "" {
		args = append(args, oImageFormat, C.CString(g.imageFormat))
	}
	if g.interlaced {
		args = append(args, oInterlaced)
	}
	if g.base != 0 {
		args = append(args, oBase, utoc(g.base))
	}
	if g.watermark != "" {
		args = append(args, oWatermark, C.CString(g.watermark))
	}
	if g.daemon != "" {
		args = append(args, oDaemon, C.CString(g.daemon))
	}
	return append(args, makeArgs(g.args)...)
}

func (e *Exporter) makeArgs(start, end time.Time, step time.Duration) []*C.char {
	args := []*C.char{
		xport,
		oStart, i64toc(start.Unix()),
		oEnd, i64toc(end.Unix()),
		oStep, i64toc(int64(step.Seconds())),
	}
	if e.maxRows != 0 {
		args = append(args, oMaxRows, utoc(e.maxRows))
	}
	if e.daemon != "" {
		args = append(args, oDaemon, C.CString(e.daemon))
	}
	return append(args, makeArgs(e.args)...)
}

func parseInfoKey(ik string) (kname, kkey string, kid int) {
	kid = -1
	o := strings.IndexRune(ik, '[')
	if o == -1 {
		kname = ik
		return
	}
	c := strings.IndexRune(ik[o+1:], ']')
	if c == -1 {
		kname = ik
		return
	}
	c += o + 1
	kname = ik[:o] + ik[c+1:]
	kkey = ik[o+1 : c]
	if id, err := strconv.Atoi(kkey); err == nil && id >= 0 {
		kid = id
	}
	return
}

func updateInfoValue(i *C.struct_rrd_info_t, v interface{}) interface{} {
	switch i._type {
	case C.RD_I_VAL:
		return float64(*(*C.rrd_value_t)(unsafe.Pointer(&i.value[0])))
	case C.RD_I_CNT:
		return uint(*(*C.ulong)(unsafe.Pointer(&i.value[0])))
	case C.RD_I_STR:
		return C.GoString(*(**C.char)(unsafe.Pointer(&i.value[0])))
	case C.RD_I_INT:
		return int(*(*C.int)(unsafe.Pointer(&i.value[0])))
	case C.RD_I_BLO:
		blob := *(*C.rrd_blob_t)(unsafe.Pointer(&i.value[0]))
		b := C.GoBytes(unsafe.Pointer(blob.ptr), C.int(blob.size))
		if v == nil {
			return b
		}
		return append(v.([]byte), b...)
	}

	return nil
}

func parseRRDInfo(i *C.rrd_info_t) map[string]interface{} {
	defer C.rrd_info_free(i)

	r := make(map[string]interface{})
	for w := (*C.struct_rrd_info_t)(i); w != nil; w = w.next {
		kname, kkey, kid := parseInfoKey(C.GoString(w.key))
		v, ok := r[kname]
		switch {
		case kid != -1:
			var a []interface{}
			if ok {
				a = v.([]interface{})
			}
			if len(a) < kid+1 {
				oldA := a
				a = make([]interface{}, kid+1)
				copy(a, oldA)
			}
			a[kid] = updateInfoValue(w, a[kid])
			v = a
		case kkey != "":
			var m map[string]interface{}
			if ok {
				m = v.(map[string]interface{})
			} else {
				m = make(map[string]interface{})
			}
			old, _ := m[kkey]
			m[kkey] = updateInfoValue(w, old)
			v = m
		default:
			v = updateInfoValue(w, v)
		}
		r[kname] = v
	}
	return r
}

func parseGraphInfo(i *C.rrd_info_t) (gi GraphInfo, img []byte) {
	inf := parseRRDInfo(i)
	if v, ok := inf["image_info"]; ok {
		gi.Print = append(gi.Print, v.(string))
	}
	for k, v := range inf {
		if k == "print" {
			for _, line := range v.([]interface{}) {
				gi.Print = append(gi.Print, line.(string))
			}
		}
	}
	if v, ok := inf["image_width"]; ok {
		gi.Width = v.(uint)
	}
	if v, ok := inf["image_height"]; ok {
		gi.Height = v.(uint)
	}
	if v, ok := inf["value_min"]; ok {
		gi.Ymin = v.(float64)
	}
	if v, ok := inf["value_max"]; ok {
		gi.Ymax = v.(float64)
	}
	if v, ok := inf["image"]; ok {
		img = v.([]byte)
	}
	return
}

func (g *Grapher) graph(filename string, start, end time.Time) (GraphInfo, []byte, error) {
	var i *C.rrd_info_t
	args := g.makeArgs(filename, start, end)

	mutex.Lock() // rrd_graph_v isn't thread safe
	defer mutex.Unlock()

	err := makeError(C.rrdGraph(
		&i,
		C.int(len(args)),
		&args[0],
	))

	if err != nil {
		return GraphInfo{}, nil, err
	}
	gi, img := parseGraphInfo(i)

	return gi, img, nil
}

// Info returns information about RRD file.
func Info(filename string) (map[string]interface{}, error) {
	fn := C.CString(filename)
	defer freeCString(fn)
	var i *C.rrd_info_t
	err := makeError(C.rrdInfo(&i, fn))
	if err != nil {
		return nil, err
	}
	return parseRRDInfo(i), nil
}

// Fetch retrieves data from RRD file.
func Fetch(filename, cf string, start, end time.Time, step time.Duration) (FetchResult, error) {
	fn := C.CString(filename)
	defer freeCString(fn)
	cCf := C.CString(cf)
	defer freeCString(cCf)
	cStart := C.time_t(start.Unix())
	cEnd := C.time_t(end.Unix())
	cStep := C.ulong(step.Seconds())
	var (
		ret      C.int
		cDsCnt   C.ulong
		cDsNames **C.char
		cData    *C.double
	)
	err := makeError(C.rrdFetch(&ret, fn, cCf, &cStart, &cEnd, &cStep, &cDsCnt, &cDsNames, &cData))
	if err != nil {
		return FetchResult{filename, cf, start, end, step, nil, 0, nil}, err
	}

	start = time.Unix(int64(cStart), 0)
	end = time.Unix(int64(cEnd), 0)
	step = time.Duration(cStep) * time.Second
	dsCnt := int(cDsCnt)

	dsNames := make([]string, dsCnt)
	for i := 0; i < dsCnt; i++ {
		dsName := C.arrayGetCString(cDsNames, C.int(i))
		dsNames[i] = C.GoString(dsName)
		C.free(unsafe.Pointer(dsName))
	}
	C.free(unsafe.Pointer(cDsNames))

	rowCnt := (int(cEnd)-int(cStart))/int(cStep) + 1
	valuesLen := dsCnt * rowCnt
	values := make([]float64, valuesLen)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&values)))
	sliceHeader.Cap = valuesLen
	sliceHeader.Len = valuesLen
	sliceHeader.Data = uintptr(unsafe.Pointer(cData))
	return FetchResult{filename, cf, start, end, step, dsNames, rowCnt, values}, nil
}

// FreeValues free values memory allocated by C.
func (r *FetchResult) FreeValues() {
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&r.values)))
	C.free(unsafe.Pointer(sliceHeader.Data))
}

// Export data from RRD file(s)
func (e *Exporter) xport(start, end time.Time, step time.Duration) (XportResult, error) {
	cStart := C.time_t(start.Unix())
	cEnd := C.time_t(end.Unix())
	cStep := C.ulong(step.Seconds())
	args := e.makeArgs(start, end, step)

	mutex.Lock()
	defer mutex.Unlock()

	var (
		ret      C.int
		cXSize   C.int
		cColCnt  C.ulong
		cLegends **C.char
		cData    *C.double
	)
	err := makeError(C.rrdXport(
		&ret,
		C.int(len(args)),
		&args[0],
		&cXSize, &cStart, &cEnd, &cStep, &cColCnt, &cLegends, &cData,
	))
	if err != nil {
		return XportResult{start, end, step, nil, 0, nil}, err
	}

	start = time.Unix(int64(cStart), 0)
	end = time.Unix(int64(cEnd), 0)
	step = time.Duration(cStep) * time.Second
	colCnt := int(cColCnt)

	legends := make([]string, colCnt)
	for i := 0; i < colCnt; i++ {
		legend := C.arrayGetCString(cLegends, C.int(i))
		legends[i] = C.GoString(legend)
		C.free(unsafe.Pointer(legend))
	}
	C.free(unsafe.Pointer(cLegends))

	rowCnt := (int(cEnd)-int(cStart))/int(cStep) + 1
	valuesLen := colCnt * rowCnt
	values := make([]float64, valuesLen)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&values)))
	sliceHeader.Cap = valuesLen
	sliceHeader.Len = valuesLen
	sliceHeader.Data = uintptr(unsafe.Pointer(cData))
	return XportResult{start, end, step, legends, rowCnt, values}, nil
}

// FreeValues free values memory allocated by C.
func (r *XportResult) FreeValues() {
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&r.values)))
	C.free(unsafe.Pointer(sliceHeader.Data))
}
