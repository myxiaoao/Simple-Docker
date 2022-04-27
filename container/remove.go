package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"simple-docker/common"
)

// RemoveContainer 删除容器
// 这里需要明确，我们是不能删除正在运行的容器的，也就是 status 为 Running 的容器，
// 所以我们删除容器之前一定要先停止容器，然后我们将该容器生成的一系列文件夹删除即可
func RemoveContainer(containerName string) {
	info, err := getContainerInfo(containerName)
	if err != nil {
		logrus.Errorf("get container info, err: %v", err)
		return
	}
	// 只能删除停止状态的容器
	if info.Status != common.Stop {
		logrus.Errorf("can't remove running container")
		return
	}
	dir := path.Join(common.DefaultContainerInfoPath, containerName)
	err = os.RemoveAll(dir)
	if err != nil {
		logrus.Errorf("remove container dir: %s, err: %v", dir, err)
		return
	}
}
