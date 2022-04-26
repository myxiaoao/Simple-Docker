package sub_system

// 内存限制实例

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type MemorySubSystem struct {
}

func (*MemorySubSystem) Name() string {
	return "memory"
}

func (m *MemorySubSystem) Set(CGroupPath string, res *ResourceConfig) error {
	subSystemCGroupPath, err := GetCGroupPath(m.Name(), CGroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", CGroupPath, err)
		return err
	}
	if res.MemoryLimit != "" {
		// 设置 CGroup 内存限制，
		// 将这个限制写入到 CGroup 对应目录的 memory.limit_in_bytes 文件中即可
		err := ioutil.WriteFile(path.Join(subSystemCGroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemorySubSystem) Remove(CGroupPath string) error {
	subSystemCGroupPath, err := GetCGroupPath(m.Name(), CGroupPath, true)
	if err != nil {
		return err
	}
	return os.RemoveAll(subSystemCGroupPath)
}

func (m *MemorySubSystem) Apply(CGroupPath string, pid int) error {
	subSystemCGroupPath, err := GetCGroupPath(m.Name(), CGroupPath, true)
	if err != nil {
		return err
	}
	tasksPath := path.Join(subSystemCGroupPath, "tasks")
	err = ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		logrus.Errorf("write pid to tasks, path: %s, pid: %d, err: %v", tasksPath, pid, err)
		return err
	}
	return nil
}
