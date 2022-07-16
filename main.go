package main

import (
	"S3ObjectStorageFileBrowser/Object"
	"S3ObjectStorageFileBrowser/Object/Config"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	var pathConfig string

	app := &cli.App{
		Name:     "S3 Object Storage File WEB Browser",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "xinjiajuan",
				Email: "itpours@qq.com",
			},
		},
		Copyright: "(c) 2022 XINJIAJUAN",
		Usage:     "基于S3对象储存的网络浏览器",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       "config.yaml",
				Usage:       "指定YAML配置文件路径",
				Aliases:     []string{"c"},
				Required:    true,
				Destination: &pathConfig,
			}, /*
				&cli.StringFlag{
					Name:     "test",
					Value:    "",
					Usage:    "",
					Aliases:  []string{"t"},
					Required: true,
					//Destination: &pathConfig,
				},*/
		},
		Action: func(cCtx *cli.Context) error {
			//判断文件夹是否存在
			_, err := os.Stat(cCtx.String("config")) //os.Stat获取文件信息
			if err != nil {
				return cli.Exit("配置文件不存在!", 86)
			} else {
				config := Config.ReadConfig(cCtx.String("config"))
				Object.MakeClientObject(config)
			}
			return nil
		},
		/*
			Commands: []*cli.Command{
				{Name: "config",
					Aliases: []string{"c"},
					Usage:   "指定YAML配置文件路径",
					Action: func(*cli.Context) error {
						println("23456")
						return nil
					},
				},
			},*/
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
