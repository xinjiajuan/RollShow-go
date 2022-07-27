package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"rollshow/Object"
	"rollshow/Object/Config"
	"time"
)

func main() {
	var pathConfig string

	app := &cli.App{
		Name:     Config.AppName,
		Version:  Config.Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "xinjiajuan",
				Email: "itpours@qq.com",
			},
		},
		Copyright: "(c) 2022 xinjiajuan",
		Usage:     Config.Usage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "config",
				//Value:       "config.yaml",
				Usage:       "指定YAML配置文件路径",
				Aliases:     []string{"c"},
				Required:    true,
				Destination: &pathConfig,
			},
		},
		Action: func(cCtx *cli.Context) error {
			//判断文件夹是否存在
			_, err := os.Stat(cCtx.String("config")) //os.Stat获取文件信息
			if err != nil {
				return cli.Exit("配置文件不存在!", 86)
			} else {
				config := Config.ReadConfig(cCtx.String("config"))
				Object.MakeS3HttpServer(config)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
