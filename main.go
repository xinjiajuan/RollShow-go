/*
Copyright <2022> 新加卷

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
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
