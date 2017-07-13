// +build cgo
// +build !appengine

package base

import "runtime"

func numCgoCall() int64 {
	return runtime.NumCgoCall()
}
