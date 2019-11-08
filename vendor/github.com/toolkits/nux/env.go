package nux

import (
	"fmt"
	"os"
	"strings"
)

const nuxRootFs = "NUX_ROOTFS"

func init() {
	fmt.Printf("use %s = %s as file system root", nuxRootFs, Root())
}

// Root 获取系统变量
func Root() string {
	root := os.Getenv(nuxRootFs)
	if !strings.HasPrefix(root, string(os.PathSeparator)) {
		return ""
	}
	root = strings.TrimSuffix(root, string(os.PathSeparator))
	if pathExists(root) {
		return root
	}
	return ""
}

func pathExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir()
	}
	return false
}
