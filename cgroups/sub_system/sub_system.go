package sub_system

// 资源限制接口

// ResourceConfig 资源限制配置
type ResourceConfig struct {
	// 内存限制
	MemoryLimit string
	// Cpu 时间片权重
	CpuShare string
	// Cpu 核数
	CpuSet string
}

// SubSystem 将 CGroup 抽象成 path，因为 hierarchy 中，CGroup 便是虚拟的路径地址
type SubSystem interface {
	// Name 返回 sub system 名字，如 cpu，memory
	Name() string
	// Set 设置 CGroup 在这个 sub system 中的资源限制
	Set(CGroupPath string, res *ResourceConfig) error
	// Remove 移除这个 CGroup 限制
	Remove(CGroupPath string) error
	// Apply 将某个进程添加到 CGroup 中，即将进程 ID 添加到 tasks 中。
	Apply(CGroupPath string, pid int) error
}

var (
	SubSystems = []SubSystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpuSetSubSystem{},
	}
)
