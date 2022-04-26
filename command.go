package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"simple-docker/cgroups/sub_system"
	"simple-docker/container"
)

// 创建 namespace 隔离容器进程
// 启动容器
var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container with namespace and CGroups limit",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "docker volume",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "detach container",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
	},
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing contaniner")
		}

		res := &sub_system.ResourceConfig{
			MemoryLimit: ctx.String("m"),
			CpuSet:      ctx.String("cpuset"),
			CpuShare:    ctx.String("cpushare"),
		}

		// cmdArray 为容器运行后，执行的第一个命令信息
		// cmdArray[0] 为命令内容, 后面的为命令参数
		var cmdArray []string
		for _, arg := range ctx.Args() {
			cmdArray = append(cmdArray, arg)
		}

		tty := ctx.Bool("ti")
		volume := ctx.String("v")
		detach := ctx.Bool("d")

		if tty && detach {
			return fmt.Errorf("ti and d paramter can not both provided")
		}

		containerName := ctx.String("name")
		Run(cmdArray, tty, res, volume, containerName)
		return nil
	},
}

// 初始化容器内容，挂载 proc 文件系统，运行用户执行程序
var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(ctx *cli.Context) error {
		logrus.Infof("Init come on.")
		return container.RunContainerInitProcess()
	},
}

// 日志命令行
var logCommand = cli.Command{
	Name:  "logs",
	Usage: "look container log",
	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container name")
		}
		containerName := ctx.Args().Get(0)
		container.LookContainerLog(containerName)
		return nil
	},
}
