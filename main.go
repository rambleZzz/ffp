package main

import (
	"github.com/rambleZzz/ffp/core"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	var (
		url     string
		file    string
		threads int
		o2Excel bool
	)

	var app = &cli.App{
		Name:    "FFP - FastFingerPrint",
		Version: "v1.0",
		Usage:   "一款基于GO语言编写的支持调用httpx、observerWard指纹识别结果、cdn检测、无cdn ip及c段去重统计、ip归属地查询于一体的快速自动化资产指纹识别工具.",
		UsageText: "ffp -u example.com (单一目标扫描)\n" +
			"ffp -f target.txt 	(批量目标扫描)\n" +
			"\n" +
			"-noe表示不输出到excel，此参数默认为false,不加此参数会输出,如果-noe不会输出至excel\n" +
			"-t参数为线程默认为5\n" +
			"-u 或-f参数为必填项(-u和-f不可同时)，其他参数为可选项，更多参数请参考GLOBAL OPTIONS",
		Action: core.TaskNew,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "从txt文件中读取内容 如:-f target.txt",
				Destination: &file,
			},
			&cli.StringFlag{
				Name:        "url",
				Aliases:     []string{"u"},
				Usage:       "target参数为需要检测的目标 如:-u example.com",
				Destination: &url,
			},
			&cli.IntFlag{
				Name:        "threads",
				Aliases:     []string{"t"},
				Value:       5,
				Usage:       "扫描线程数，如：-t 5",
				Destination: &threads,
			},
			&cli.BoolFlag{
				Name:        "notO2Excel",
				Aliases:     []string{"noe"},
				Value:       false, //false代表默认输出
				Usage:       "是否不输出到excel, 如：-noe 不输出",
				Destination: &o2Excel,
			},
			// 省略剩余的 StringFlag...
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
