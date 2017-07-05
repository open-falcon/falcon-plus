package str

import (
	"crypto/md5"
	"fmt"
)

func Md5Encode(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}
