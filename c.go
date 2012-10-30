package rrd

/*
#include <stdlib.h>
#include <rrd.h>
*/
import "C"
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

func (c *Create) create() error {
	filename := C.CString(c.filename)
	defer freeCString(filename)
	a, n := makeArgs(args)
	defer freeArgs(a)
	//C.rrd_create_r(c.filename, C.
}

/*fmt.Sprint(start.Unix())
a[4], a[5] = "--step", fmt.Sprint((int64(step)+0.5e9)/1e9)*/
