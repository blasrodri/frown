package lsof

import (
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func GetUserProcessList() ([]*Process, error) {
	listUserPids, err := getUserPids()
	listProcesses := make([]*Process, len(listUserPids))

	if err != nil {
		return nil, err
	}

	for i, pid := range listUserPids {
		process, err := newProcess(pid)
		if err != nil {
			continue
		}
		listProcesses[i] = process
	}

	return listProcesses, nil
}

func getUserPids() ([]int, error) {
	userInfo, _ := user.Current()
	userUID, _ := strconv.Atoi(userInfo.Uid)
	userGUID, err := strconv.Atoi(userInfo.Uid)
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
		// 1. Ensure that there is a dir with the fd name
		// 2. Ensure that the user has permissions on that dir
		// 3. Ensure that the user has permissions on the status file
		info, err := os.Stat(cwd + "/" + dir.Name())
		if err != nil {
			// Assume that the pid does not exist anymore
			// Skip it
			continue
		}
		var UID int
		var GID int
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			UID = int(stat.Uid)
			GID = int(stat.Gid)
			if UID == userUID || GID == userGUID {
				pid, err := strconv.Atoi(info.Name())
				if err != nil {
					return nil, err
				}
				infoStatusFd, err := os.Stat(cwd + "/" + dir.Name() + "/" + "status")
				if err != nil {
					return nil, err
				}
				if stat, ok := infoStatusFd.Sys().(*syscall.Stat_t); ok {
					UID = int(stat.Uid)
					GID = int(stat.Gid)

					if UID == userUID || GID == userGUID {
						pids = append(pids, pid)
					}
				}
			}
		}
	}
	return pids, nil
}
