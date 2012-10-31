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

func (g *Grapher) SaveGraph(filename string, start, end time.Time) (GraphInfo, error) {
	var info *C.rrd_info_t
	args := g.makeArgs(filename, start, end)

	g.m.Lock() // rrd_graph_v isn't thread safe
	err := makeError(C.rrdGraph(
		&info,
		C.int(len(args)),
		&args[0],
	))
	g.m.Unlock()

	if err != nil {
		return GraphInfo{}, err
	}
	return GraphInfo{}, nil
}
