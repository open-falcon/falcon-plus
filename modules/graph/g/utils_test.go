package g

import (
	"fmt"
	"testing"
)

func Test_RrdFileName(t *testing.T) {
	if RrdFileName("/basedir", "b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) !=
		RrdFileName_orig("/basedir", "b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) {
		t.Error("not match with orig func")
	}

	if RrdFileName("/basedir", "b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) !=
		"/basedir/b0/b026324c6904b2a9cb4b88d6d61c81d1_GAUGE_10.rrd" {
		t.Error("not match")
	}
}

func Test_FormRrdCacheKey(t *testing.T) {
	if FormRrdCacheKey("b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) !=
		FormRrdCacheKey_orig("b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) {
		t.Error("not match")
	}

	if FormRrdCacheKey("b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10) !=
		"b026324c6904b2a9cb4b88d6d61c81d1_GAUGE_10" {
		t.Error("not match")
	}
}

func Benchmark_RrdFileName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RrdFileName("/basedir", "b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10)
	}
}

func Benchmark_RrdFileName_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RrdFileName_orig("/basedir", "b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10)
	}
}

func Benchmark_FormRrdCacheKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormRrdCacheKey("b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10)
	}
}

func Benchmark_FormRrdCacheKey_orig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormRrdCacheKey_orig("b026324c6904b2a9cb4b88d6d61c81d1", "GAUGE", 10)
	}
}

func RrdFileName_orig(baseDir string, md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s/%s/%s_%s_%d.rrd", baseDir, md5[0:2], md5, dsType, step)
}

func FormRrdCacheKey_orig(md5 string, dsType string, step int) string {
	return fmt.Sprintf("%s_%s_%d", md5, dsType, step)
}
