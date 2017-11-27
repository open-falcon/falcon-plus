package rrdtool

import (
	"strconv"

	"github.com/open-falcon/falcon-plus/modules/graph/g"
)

// getIndex
// 输入: md5
// 输出: slice的index
func getIndex(md5 string) (index int) {
	batchNum := g.Config().IOWorkerNum
	firstBytesSize := g.Config().FirstBytesSize

	if batchNum <= 1 || len(md5) < firstBytesSize || firstBytesSize == 0 {
		return 0
	}

	m, err := strconv.ParseInt(md5[0:firstBytesSize], 16, 64)
	if err != nil {
		return 0
	}

	return int(m) % int(batchNum)
}
