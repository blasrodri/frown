package lsof

import (
	"io/ioutil"
	"os"
	"os/user"
	"syscall"
	"strconv"
)

func GetUserPids() ([]int, error) {
	userInfo, err := user.Current()
	userUid, err:= strconv.Atoi(userInfo.Uid)
	userGuid, err:= strconv.Atoi(userInfo.Uid)
	if err != nil {
		return nil, err
	}

	cwd := "/proc"
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0)
	for _, dir := range files {
		if !dir.IsDir() {
			continue
		}
		info, _ := os.Stat(cwd + "/" + dir.Name())
		var UID int
		var GID int
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			UID = int(stat.Uid)
			GID = int(stat.Gid)
			if UID ==  userUid || GID == userGuid {
				pid, err := strconv.Atoi(info.Name())
				if err != nil {
					return nil, err
				}
				pids = append(pids,pid)
			}
		}
	}
	return pids, nil
}
