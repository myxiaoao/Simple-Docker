package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"simple-docker/common"
)

// LookContainerLog 查看容器内日志信息
func LookContainerLog(containerName string) {
	logFileName := path.Join(common.DefaultContainerInfoPath, containerName, common.ContainerLogFileName)
	file, err := os.Open(logFileName)
	if err != nil {
		logrus.Errorf("open log file, path: %s, err: %v", logFileName, err)
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("read log file, err: %v", err)
	}
	_, _ = fmt.Fprint(os.Stdout, string(bs))
}
