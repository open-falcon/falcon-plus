package g

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/toolkits/file"
)

// RRDTOOL UTILS
// 监控数据对应的rrd文件名称
func RrdFileName(baseDir string, md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s/%s/%s_%s_%d.rrd", baseDir, md5[0:2], md5, dsType, step)
}

// rrd文件是否存在
func IsRrdFileExist(filename string) bool {
	return file.IsExist(filename)
}

// 生成rrd缓存数据的key
func FormRrdCacheKey(md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s_%s_%d", md5, dsType, step)
}
func SplitRrdCacheKey(ckey string) (md5 string, dsType string, step int, err error) {
	ckey_slice := strings.Split(ckey, "_")
	if len(ckey_slice) != 3 {
		err = fmt.Errorf("bad rrd cache key: %s", ckey)
		return
	}

	md5 = ckey_slice[0]
	dsType = ckey_slice[1]
	stepInt64, err := strconv.ParseInt(ckey_slice[2], 10, 32)
	if err != nil {
		return
	}
	step = int(stepInt64)

	// return
	err = nil
	return
}
