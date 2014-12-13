package plugins

import (
	"github.com/open-falcon/agent/g"
	"github.com/toolkits/file"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	path_plugin_version = make(map[string]string)
	path_plugins        = make(map[string]map[string]*Plugin)
)

func PluginsUnder(relativePath, pluginVersion string) map[string]*Plugin {

	ver, ver_exists := path_plugin_version[relativePath]
	plugins, plugins_exists := path_plugins[relativePath]
	if ver_exists && plugins_exists && pluginVersion == ver {
		return plugins
	}

	plugins = ListPlugins(relativePath)
	path_plugins[relativePath] = plugins
	path_plugin_version[relativePath] = pluginVersion

	return plugins
}

func ListPlugins(relativePath string) map[string]*Plugin {
	ret := make(map[string]*Plugin)
	if relativePath == "" {
		return ret
	}

	dir := filepath.Join(g.Config().Plugin.Dir, relativePath)

	if !file.IsExist(dir) || file.IsFile(dir) {
		return ret
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println("can not list files under", dir)
		return ret
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		arr := strings.Split(filename, "_")
		if len(arr) < 2 {
			continue
		}

		// filename should be: $cycle_$xx
		var cycle int
		cycle, err = strconv.Atoi(arr[0])
		if err != nil {
			continue
		}

		fpath := filepath.Join(relativePath, filename)
		plugin := &Plugin{FilePath: fpath, MTime: f.ModTime().Unix(), Cycle: cycle}
		ret[fpath] = plugin
	}

	return ret
}
