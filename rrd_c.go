package rrd

/*
#include <stdlib.h>
#include "rrdfunc.h"
#cgo LDFLAGS: -lrrd_th
*/
import "C"
import "time"
import "unsafe"

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
	return Error(C.GoString(e))
}

func (c *Create) create() error {
	filename := C.CString(c.filename)
	defer freeCString(filename)
	args := makeArgs(c.args)
	defer freeArgs(args)

	e := C.rrdCreate(
		filename,
		C.ulong((c.step+time.Second/2)/time.Second),
		C.time_t(c.start.Unix()),
		C.int(len(args)),
		&args[0],
	)
	return makeError(e)
}

/*fmt.Sprint(start.Unix())
a[4], a[5] = "--step", fmt.Sprint((int64(step)+0.5e9)/1e9)*/
