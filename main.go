package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const usage = `simple-docker`

func main() {
	// 创建新的 cli 命令
	app := cli.NewApp()
	app.Name = "simple-docker"
	app.Usage = usage

	// 加载 run ,inti 命令
	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}

	// 启动前配置
	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	// 启动命令
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
