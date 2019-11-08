package nux

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"strings"
	"syscall"

	"github.com/toolkits/file"
)

// return: [][$fsSpec, $fsFile, $fsVfstype]
func ListMountPoint() ([][3]string, error) {
	contents, err := ioutil.ReadFile(Root() + "/proc/mounts")
	if err != nil {
		return nil, err
	}

	ret := make([][3]string, 0)

	reader := bufio.NewReader(bytes.NewBuffer(contents))
	for {
		line, err := file.ReadLine(reader)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return nil, err
		}

		fields := strings.Fields(string(line))
		// Docs come from the fstab(5)
		// fsSpec     # Mounted block special device or remote filesystem e.g. /dev/sda1
		// fsFile     # Mount point e.g. /data
		// fsVfstype  # File system type e.g. ext4
		// fs_mntops   # Mount options
		// fs_freq     # Dump(8) utility flags
		// fs_passno   # Order in which filesystem checks are done at reboot time

		fsSpec := fields[0]
		fsFile := fields[1]
		fsVfstype := fields[2]

		if _, exist := FSSPEC_IGNORE[fsSpec]; exist {
			continue
		}

		if _, exist := FSTYPE_IGNORE[fsVfstype]; exist {
			continue
		}

		if strings.HasPrefix(fsVfstype, "fuse") {
			continue
		}

		if IgnoreFsFile(fsFile) {
			continue
		}

		// keep /dev/xxx device with shorter fsFile (remove mount binds)
		if strings.HasPrefix(fsSpec, Root()+"/dev") {
			deviceFound := false
			for idx := range ret {
				if ret[idx][0] == fsSpec {
					deviceFound = true
					if len(fsFile) < len(ret[idx][1]) {
						ret[idx][1] = fsFile
					}
					break
				}
			}
			if !deviceFound {
				ret = append(ret, [3]string{fsSpec, fsFile, fsVfstype})
			}
		} else {
			ret = append(ret, [3]string{fsSpec, fsFile, fsVfstype})
		}
	}
	return ret, nil
}

func BuildDeviceUsage(_fsSpec, _fsFile, _fsVfstype string) (*DeviceUsage, error) {
	ret := &DeviceUsage{FsSpec: _fsSpec, FsFile: _fsFile, FsVfstype: _fsVfstype}

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(_fsFile, &fs)
	if err != nil {
		return nil, err
	}

	// blocks
	used := fs.Blocks - fs.Bfree
	ret.BlocksAll = uint64(fs.Frsize) * fs.Blocks
	ret.BlocksUsed = uint64(fs.Frsize) * used
	ret.BlocksFree = uint64(fs.Frsize) * fs.Bavail
	if fs.Blocks == 0 {
		ret.BlocksUsedPercent = 0
		ret.BlocksFreePercent = 0
	} else {
		ret.BlocksUsedPercent = float64(used) * 100.0 / float64(used+fs.Bavail)
		ret.BlocksFreePercent = 100.0 - ret.BlocksUsedPercent
	}

	// inodes
	ret.InodesAll = fs.Files
	if fs.Ffree == math.MaxUint64 {
		ret.InodesFree = 0
		ret.InodesUsed = 0
	} else {
		ret.InodesFree = fs.Ffree
		ret.InodesUsed = fs.Files - fs.Ffree
	}
	if fs.Files == 0 {
		ret.InodesUsedPercent = 0
		ret.InodesFreePercent = 0
	} else {
		ret.InodesUsedPercent = float64(ret.InodesUsed) * 100.0 / float64(ret.InodesAll)
		ret.InodesFreePercent = 100.0 - ret.InodesUsedPercent
	}

	return ret, nil
}
