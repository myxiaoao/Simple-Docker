package sub_system

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

// GetCGroupPath 获取 CGroup 在文件系统中的绝对路径
func GetCGroupPath(subSystem string, CGroupPath string, autoCreate bool) (string, error) {
	CGroupRootPath, err := findCGroupMountPoint(subSystem)
	if err != nil {
		logrus.Errorf("find cgroup mount point, err: %s", err.Error())
		return "", err
	}
	CGroupTotalPath := path.Join(CGroupRootPath, CGroupPath)
	_, err = os.Stat(CGroupTotalPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(CGroupTotalPath, 0755); err != nil {
			return "", err
		}
	}

	return CGroupTotalPath, nil
}

// 找到挂载了 subsystem 的 hierarchy CGroup 根节点所在的目录
func findCGroupMountPoint(subSystem string) (string, error) {
	// Linux 系统的 /proc/self/mountinfo 记录当前系统所有挂载文件系统的信息
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Errorf("mountinfo file close failed, err %v", err)
		}
	}(f)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subSystem && len(fields) > 4 {
				return fields[4], nil
			}
		}
	}

	return "", scanner.Err()
}
