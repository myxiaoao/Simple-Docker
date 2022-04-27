package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"simple-docker/cgroups"
	"simple-docker/cgroups/sub_system"
	"simple-docker/common"
	"simple-docker/container"
	"strings"
)

// Run 命令主要就是启动一个容器，然后对该进程设置隔离
func Run(cmdArray []string, tty bool, res *sub_system.ResourceConfig, containerName, imageName, volume string, envs []string) {
	id := container.GenContainerID(10)
	if containerName == "" {
		containerName = id
	}
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envs)
	if parent == nil {
		logrus.Errorf("failed to new parent process")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start failed, err: %v", err)
		return
	}

	// 记录容器信息
	err := container.RecordContainerInfo(parent.Process.Pid, cmdArray, containerName, id)
	if err != nil {
		logrus.Errorf("record container info, err: %v", err)
	}

	// 添加资源限制
	CGroupManager := cgroups.NewCGroupManager("simple-docker")
	// 删除资源限制
	defer CGroupManager.Destroy()
	// 设置资源限制
	CGroupManager.Set(res)
	// 将容器进程，加入到各个 subsystem 挂载对应的 CGroup 中
	CGroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmdArray, writePipe)

	if tty {
		// 等待父进程结束
		err := parent.Wait()
		if err != nil {
			logrus.Errorf("parent wait, err: %v", err)
		}
		// 删除容器工作空间
		err = container.DeleteWorkSpace(common.RootPath, volume)
		if err != nil {
			logrus.Errorf("delete work space, err: %v", err)
		}
		// 删除容器信息
		container.DeleteContainerInfo(containerName)
	}
	if err != nil {
		return
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	_, _ = writePipe.WriteString(command)
	_ = writePipe.Close()
}
