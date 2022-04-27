package common

const (
	RootPath   = "/root/"
	MntPath    = "/root/mnt/"
	WriteLayer = "writeLayer"
)

const (
	DefaultContainerInfoPath = "/var/run/simple-docker/"
	ContainerInfoFileName    = "config.json"
	ContainerLogFileName     = "container.log"
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
