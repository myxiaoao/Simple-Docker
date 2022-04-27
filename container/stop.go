package container

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"simple-docker/common"
	"strconv"
	"syscall"
)

// StopContainer 停止容器，修改容器状态
// 我们通过 config.json 记录了容器的基本信息，
// 其中就有一个 status 字段用来记录容器的状态，和一个 PID 字段记录容器的 init 进程在宿主机上的 Pid，
// 我们停止容器，也就是将该 pid 进程杀死，并更新 status 状态即可
func StopContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info, err: %v", err)
		return
	}
	if info.Pid != "" {
		pid, _ := strconv.Atoi(info.Pid)
		// 杀死进程
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			logrus.Errorf("stop container, pid: %d, err: %v", pid, err)
			return
		}
		// 修改容器状态
		info.Status = common.Stop
		info.Pid = ""
		bs, _ := json.Marshal(info)
		fileName := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerInfoFileName)
		err := ioutil.WriteFile(fileName, bs, 0622)
		if err != nil {
			logrus.Errorf("write container config.json, err: %v", err)
		}
	}
}
