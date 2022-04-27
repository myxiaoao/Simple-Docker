package sub_system

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetCGroupPath(t *testing.T) {
	logrus.Infof(findCGroupMountPoint("memory"))
	logrus.Infof(findCGroupMountPoint("cpu"))
	logrus.Infof(findCGroupMountPoint("cpuset"))
}
