package g

import (
	nproc "github.com/toolkits/proc"
)

var (
	RecvDataTrace  = nproc.NewDataTrace("RecvDataTrace", 5)
	RecvDataFilter = nproc.NewDataFilter("RecvDataFilter", 5)
)
