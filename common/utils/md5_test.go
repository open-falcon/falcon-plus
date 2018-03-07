package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"
)

func origMd5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Test_Md5(t *testing.T) {
	if Md5("1234567890123") != origMd5("1234567890123") {
		t.Error("not expect")
	}
}

func Benchmark_Md5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Md5("1234567890123")
	}
}

func Benchmark_Md5_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		origMd5("1234567890123")
	}
}
