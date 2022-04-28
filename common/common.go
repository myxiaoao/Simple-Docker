package common

const (
	RootPath   = "/root/"
	MntPath    = "/root/mnt/"
	WriteLayer = "writeLayer"
)

const (
	Running = "running"
	Stop    = "stopped"
	Exit    = "exited"
)

const (
	EnvExecPid = "docker_pid"
	EnvExecCmd = "docker_cmd"
)

const (
	DefaultContainerInfoPath = "/var/run/simple-docker/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"
)

const (
	DefaultNetworkPath   = "/var/run/simple-docker/network/network/"
	DefaultAllocatorPath = "/var/run/simple-docker/network/ipam/subnet.json"
)
