package rrd

/*
#include <stdlib.h>
#include <rrd.h>
#include "rrdfunc.h"
#cgo LDFLAGS: -lrrd_th
*/
import "C"
import (
	"fmt"
	"math"
	"time"
	"unsafe"
)

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
	graphv           = C.CString("graphv")
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

	oLogarithmic = C.CString("-o")

	oNoLegand = C.CString("-g")

	oLazy = C.CString("-z")

	oColor = C.CString("-c")

	oSlopeMode = C.CString("-E")
)

func (g *Grapher) makeArgs(filename string, start, end time.Time) []*C.char {
	args := []*C.char{
		graphv, C.CString(filename),
		oStart, C.CString(fmt.Sprint(start.Unix())),
		oEnd, C.CString(fmt.Sprint(end.Unix())),
		oTitle, C.CString(g.title),
		oVlabel, C.CString(g.vlabel),
	}
	if g.width != 0 {
		args = append(args, oWidth, C.CString(fmt.Sprint(g.width)))
	}
	if g.height != 0 {
		args = append(args, oHeight, C.CString(fmt.Sprint(g.height)))
	}
	if g.upperLimit != -math.MaxFloat64 {
		args = append(args, oUpperLimit, C.CString(fmt.Sprint(g.upperLimit)))
	}
	if g.lowerLimit != math.MaxFloat64 {
		args = append(args, oLowerLimit, C.CString(fmt.Sprint(g.lowerLimit)))
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
	if g.noLegand {
		args = append(args, oNoLegand)
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
	return append(args, makeArgs(g.args)...)
}

func parseRRDInfo(i *C.rrd_info_t) map[string]interface{} {
	defer C.rrd_info_free(i)

	r := make(map[string]interface{})
	for w := (*C.struct_rrd_info_t)(i); w != nil; w = w.next {
		k := C.GoString(w.key)
		switch w._type {
		case C.RD_I_VAL:
			r[k] = float64(*(*C.rrd_value_t)(unsafe.Pointer(&w.value[0])))
		case C.RD_I_CNT:
			r[k] = uint(*(*C.ulong)(unsafe.Pointer(&w.value[0])))
		case C.RD_I_STR:
			s := C.GoString(*(**C.char)(unsafe.Pointer(&w.value[0])))
			r[k] = s
		case C.RD_I_INT:
			r[k] = int(*(*C.int)(unsafe.Pointer(&w.value[0])))
		case C.RD_I_BLO:
			blob := *(*C.rrd_blob_t)(unsafe.Pointer(&w.value[0]))
			b := C.GoBytes(unsafe.Pointer(blob.ptr), C.int(blob.size))
			if v, ok := r[k]; ok {
				r[k] = append(v.([]byte), b...)
			} else {
				r[k] = b
			}
		}
	}
	return r
}

func parseGraphInfo(i *C.rrd_info_t) (gi GraphInfo, img []byte) {
	inf := parseRRDInfo(i)
	if v, ok := inf["image_info"]; ok {
		gi.Print = append(gi.Print, v.(string))
	}
	for k, v := range inf {
		if k[:5] == "print" {
			gi.Print = append(gi.Print, v.(string))
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

	g.m.Lock() // rrd_graph_v isn't thread safe
	err := makeError(C.rrdGraph(
		&i,
		C.int(len(args)),
		&args[0],
	))
	g.m.Unlock()

	if err != nil {
		return GraphInfo{}, nil, err
	}
	gi, img := parseGraphInfo(i)

	return gi, img, nil
}

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
