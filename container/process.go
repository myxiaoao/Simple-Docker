package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path"
	"simple-docker/common"
	"syscall"
)

// NewParentProcess 创建一个会隔离的 namespace 进程的 Command
// 通过 /proc/self/exe init 来调用自身我们定义的 initCommand 命令，然后给该进程设置隔离信息。
func NewParentProcess(tty bool, volume, containerName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, _ := os.Pipe()
	// 调用自身，传入 init 参数，也就是执行 initCommand
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC, // Cloneflags linux 特定参数
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		// 把日志输出到文件里
		logDir := path.Join(common.DefaultContainerInfoPath, containerName)
		if _, err := os.Stat(logDir); err != nil && os.IsNotExist(err) {
			err := os.MkdirAll(logDir, os.ModePerm)
			if err != nil {
				logrus.Errorf("mkdir container log, err: %v", err)
			}
		}
		logFileName := path.Join(logDir, common.ContainerLogFileName)
		file, err := os.Create(logFileName)
		if err != nil {
			logrus.Errorf("create log file, err: %v", err)
		}
		// 将 cmd 的输出流改到文件流中
		cmd.Stdout = file
	}

	// 设置额外文件句柄
	cmd.ExtraFiles = []*os.File{
		readPipe,
	}

	return cmd, writePipe
}
