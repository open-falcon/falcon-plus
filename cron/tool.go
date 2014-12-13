package cron

import (
	"bytes"
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/file"
	"os/exec"
	"strings"
)

func GetCurrPluginVersion() string {
	if !g.Config().Plugin.Enabled {
		return "plugin not enabled"
	}

	pluginDir := g.Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		return "plugin dir not existent"
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error:%s", err.Error())
	}

	return strings.TrimSpace(out.String())
}
