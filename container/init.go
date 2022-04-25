package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunContainerInitProcess 本容器执行的第一个进程
// 使用 mount 挂载 proc 文件系统
// 以便后面通过 `ps` 等系统命令查看当前进程资源情况
func RunContainerInitProcess() error {
	// 读取用户命令行
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("get user command not in run contariner")
	}

	// 挂载
	err := setUpMount()
	if err != nil {
		logrus.Errorf("set up mount, err: %v", err)
		return err
	}

	// 在系统环境 PATH 中寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("look %s path, err: %v", cmdArray[0], err)
		return err
	}

	err = syscall.Exec(path, cmdArray[0:], os.Environ())
	if err != nil {
		return err
	}

	return nil
}

// 读取用户命令行
func readUserCommand() []string {
	// 指定 index 为 3 的文件描述符
	// 也就是 cmd.ExtraFiles 中，我们传递过来的 readPipe
	pipe := os.NewFile(uintptr(3), "pipe")
	bs, err := ioutil.ReadAll(pipe)

	if err != nil {
		logrus.Errorf("read pipe, err:%v", err)
	}

	msg := string(bs)

	return strings.Split(msg, " ")
}

// 挂载
func setUpMount() error {
	// systemd 加入 linux 以后，mount namespace 就变成 shared by default， 所以必须显示声明
	// 声明新的 mount namespace 独立
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		return err
	}
	// mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("mount proc , err: %v", err)
	}

	return nil
}
