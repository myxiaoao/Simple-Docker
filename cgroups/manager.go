package cgroups

// 资源限制管理器

import (
	"github.com/sirupsen/logrus"
	"simple-docker/cgroups/sub_system"
)

type CGroupManager struct {
	Path string
}

func NewCGroupManager(path string) *CGroupManager {
	return &CGroupManager{Path: path}
}

func (c *CGroupManager) Set(res *sub_system.ResourceConfig) {
	for _, subSystem := range sub_system.SubSystems {
		err := subSystem.Set(c.Path, res)
		if err != nil {
			logrus.Errorf("set %s err %v", subSystem.Name(), err)
		}
	}
}

func (c *CGroupManager) Apply(pid int) {
	for _, subSystem := range sub_system.SubSystems {
		err := subSystem.Apply(c.Path, pid)
		if err != nil {
			logrus.Errorf("apply task, err: %v", err)
		}
	}
}

func (c *CGroupManager) Destroy() {
	for _, subSystem := range sub_system.SubSystems {
		err := subSystem.Remove(c.Path)
		if err != nil {
			logrus.Errorf("remove %s err: %v", subSystem.Name(), err)
		}
	}
}
