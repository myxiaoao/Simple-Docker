package sub_system

// Cpu 核数限制实例

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSetSubSystem struct {
	apply bool
}

func (*CpuSetSubSystem) Name() string {
	return "cpuset"
}

func (c *CpuSetSubSystem) Set(CGroupPath string, res *ResourceConfig) error {
	subSystemCGroupPath, err := GetCGroupPath(c.Name(), CGroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", CGroupPath, err)
		return err
	}
	if res.CpuSet != "" {
		c.apply = true
		err := ioutil.WriteFile(path.Join(subSystemCGroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644)
		if err != nil {
			logrus.Errorf("failed to write file cpuset.cpus, err: %+v", err)
			return err
		}
	}
	return nil
}

func (c *CpuSetSubSystem) Remove(CGroupPath string) error {
	subSystemCGroupPath, err := GetCGroupPath(c.Name(), CGroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subSystemCGroupPath)
}

func (c *CpuSetSubSystem) Apply(CGroupPath string, pid int) error {
	if c.apply {
		subSystemCGroupPath, err := GetCGroupPath(c.Name(), CGroupPath, false)
		if err != nil {
			return err
		}
		tasksPath := path.Join(subSystemCGroupPath, "tasks")
		err = ioutil.WriteFile(tasksPath, []byte(strconv.Itoa(pid)), os.ModePerm)
		if err != nil {
			logrus.Errorf("write pid to tasks, path: %s, pid: %d, err: %v", tasksPath, pid, err)
			return err
		}
	}
	return nil
}
