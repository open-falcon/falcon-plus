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
	graphv  = C.CString("graphv")
	oStart  = C.CString("-s")
	oEnd    = C.CString("-e")
	oTitle  = C.CString("-t")
	oVlabel = C.CString("-v")
	oWidth  = C.CString("-w")
	oHeight = C.CString("-h")
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
	return append(args, makeArgs(g.args)...)
}

func parseRRDInfo(i *C.rrd_info_t) (gi GraphInfo, img []byte) {
	defer C.rrd_info_free(i)

	for w := (*C.struct_rrd_info_t)(i); w != nil; w = w.next {
		if C.GoString(w.key) == "image_info" {
			gi.Print = append(
				gi.Print,
				C.GoString(*(**C.char)(unsafe.Pointer(&w.value[0]))),
			)
		}
	}
	for w := (*C.struct_rrd_info_t)(i); w != nil; w = w.next {
		switch C.GoString(w.key) {
		case "image_width":
			gi.Width = uint(*(*C.ulong)(unsafe.Pointer(&w.value[0])))
		case "image_height":
			gi.Height = uint(*(*C.ulong)(unsafe.Pointer(&w.value[0])))
		case "value_min":
			gi.Ymin = float64(*(*C.rrd_value_t)(unsafe.Pointer(&w.value[0])))
		case "value_max":
			gi.Ymax = float64(*(*C.rrd_value_t)(unsafe.Pointer(&w.value[0])))
		case "print":
			gi.Print = append(
				gi.Print,
				C.GoString(*(**C.char)(unsafe.Pointer(&w.value[0]))),
			)
		case "image":
			blob := *(*C.rrd_blob_t)(unsafe.Pointer(&w.value[0]))
			buf := C.GoBytes(unsafe.Pointer(blob.ptr), C.int(blob.size))
			img = append(img, buf...)
		}
	}

	return
}

func (g *Grapher) graph(filename string, start, end time.Time) (GraphInfo, []byte, error) {
	var i *C.rrd_info_t
	args := g.makeArgs(filename, start, end)
	fmt.Println(args)

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
	gi, img := parseRRDInfo(i)

	return gi, img, nil
}
