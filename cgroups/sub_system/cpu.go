package sub_system

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CpuSubSystem struct {
	apply bool
}

func (*CpuSubSystem) Name() string {
	return "cpu"
}

func (c *CpuSubSystem) Set(CGroupPath string, res *ResourceConfig) error {
	subSystemCGroupPath, err := GetCGroupPath(c.Name(), CGroupPath, true)
	if err != nil {
		logrus.Errorf("get %s path, err: %v", CGroupPath, err)
		return err
	}
	if res.CpuShare != "" {
		c.apply = true
		err = ioutil.WriteFile(path.Join(subSystemCGroupPath, "cpu.shares"), []byte(res.CpuShare), 0644)
		if err != nil {
			logrus.Errorf("failed to write file cpu.shares, err: %+v", err)
			return err
		}
	}

	return nil
}

func (c *CpuSubSystem) Remove(CGroupPath string) error {
	subSystemCGroupPath, err := GetCGroupPath(c.Name(), CGroupPath, false)
	if err != nil {
		return err
	}

	return os.RemoveAll(subSystemCGroupPath)
}

func (c *CpuSubSystem) Apply(CGroupPath string, pid int) error {
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
