package main

import (
	"fmt"
	"git.zc0901.com/go/god/tools/god/api/apigen"
	"git.zc0901.com/go/god/tools/god/api/dartgen"
	"git.zc0901.com/go/god/tools/god/api/docgen"
	"git.zc0901.com/go/god/tools/god/api/format"
	"git.zc0901.com/go/god/tools/god/api/gogen"
	"git.zc0901.com/go/god/tools/god/api/javagen"
	"git.zc0901.com/go/god/tools/god/api/ktgen"
	"git.zc0901.com/go/god/tools/god/api/new"
	"git.zc0901.com/go/god/tools/god/api/tsgen"
	"git.zc0901.com/go/god/tools/god/api/validate"
	"git.zc0901.com/go/god/tools/god/mysql/command"
	rpc "git.zc0901.com/go/god/tools/god/rpc/cli"
	"github.com/urfave/cli"
	"os"
	"runtime"
)

var (
	BuildTime = "20201126"
	commands  = []cli.Command{
		{
			Name:   "mysql",
			Usage:  "从数据源生成MySQL模型层代码",
			Action: command.GenCodeFromDSN,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dsn",
					Usage: `数据库连接地址，如 "root:asdfasdf@tcp(192.168.0.166:3306)/nest_label"`,
				},
				cli.StringFlag{
					Name:  "table, t",
					Usage: `表名，多表以英文逗号分隔，如 "node,tag,channel"`,
				},
				cli.StringFlag{
					Name:  "dir, d",
					Usage: "目标文件夹",
				},
				cli.BoolFlag{
					Name:  "cache, c",
					Usage: "生成带缓存的数据访问层[可选]",
				},
			},
		},
		{
			Name:  "rpc",
			Usage: "generate rpc code",
			Subcommands: []cli.Command{
				{
					Name:  "new",
					Usage: `generate rpc demo service`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:     "style",
							Required: false,
							Usage:    "the file naming format, see [https://github.com/tal-tech/go-zero/tree/master/tools/goctl/config/readme.md]",
						},
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [optional]",
						},
					},
					Action: rpc.RpcNew,
				},
				{
					Name:  "template",
					Usage: `generate proto template`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "out, o",
							Usage: "the target path of proto",
						},
					},
					Action: rpc.RpcTemplate,
				},
				{
					Name:  "proto",
					Usage: `generate rpc from proto`,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "src, s",
							Usage: "the file path of the proto source file",
						},
						cli.StringSliceFlag{
							Name:  "proto_path, I",
							Usage: `native command of protoc, specify the directory in which to search for imports. [optional]`,
						},
						cli.StringFlag{
							Name:  "dir, d",
							Usage: `the target path of the code`,
						},
						cli.StringFlag{
							Name:     "style",
							Required: false,
							Usage:    "the file naming format, see [https://github.com/tal-tech/go-zero/tree/master/tools/goctl/config/readme.md]",
						},
						cli.BoolFlag{
							Name:  "idea",
							Usage: "whether the command execution environment is from idea plugin. [optional]",
						},
					},
					Action: rpc.Rpc,
				},
			},
		},
		{
			Name:  "api",
			Usage: "generate api related files",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "o",
					Usage: "the output api file",
				},
			},
			Action: apigen.ApiCommand,
			Subcommands: []cli.Command{
				{
					Name:   "new",
					Usage:  "fast create api service",
					Action: new.NewService,
				},
				{
					Name:  "format",
					Usage: "format api files",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the format target dir",
						},
						cli.BoolFlag{
							Name:     "iu",
							Usage:    "ignore update",
							Required: false,
						},
						cli.BoolFlag{
							Name:     "stdin",
							Usage:    "use stdin to input api doc content, press \"ctrl + d\" to send EOF",
							Required: false,
						},
					},
					Action: format.GoFormatApi,
				},
				{
					Name:  "validate",
					Usage: "validate api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "api",
							Usage: "validate target api file",
						},
					},
					Action: validate.GoValidateApi,
				},
				{
					Name:  "doc",
					Usage: "generate doc files",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
					},
					Action: docgen.DocCommand,
				},
				{
					Name:  "go",
					Usage: "generate go files for provided api in yaml file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
						cli.StringFlag{
							Name:     "style",
							Required: false,
							Usage:    "the file naming format, see [https://github.com/tal-tech/go-zero/tree/master/tools/goctl/config/readme.md]",
						},
					},
					Action: gogen.GoCommand,
				},
				{
					Name:  "java",
					Usage: "generate java files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: javagen.JavaCommand,
				},
				{
					Name:  "ts",
					Usage: "generate ts files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
						cli.StringFlag{
							Name:     "webapi",
							Usage:    "the web api file path",
							Required: false,
						},
						cli.StringFlag{
							Name:     "caller",
							Usage:    "the web api caller",
							Required: false,
						},
						cli.BoolFlag{
							Name:     "unwrap",
							Usage:    "unwrap the webapi caller for import",
							Required: false,
						},
					},
					Action: tsgen.TsCommand,
				},
				{
					Name:  "dart",
					Usage: "generate dart files for provided api in api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target dir",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
					},
					Action: dartgen.DartCommand,
				},
				{
					Name:  "kt",
					Usage: "generate kotlin code for provided api file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "dir",
							Usage: "the target directory",
						},
						cli.StringFlag{
							Name:  "api",
							Usage: "the api file",
						},
						cli.StringFlag{
							Name:  "pkg",
							Usage: "define package name for kotlin file",
						},
					},
					Action: ktgen.KtCommand,
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "God 代码生成器"
	app.Version = fmt.Sprintf("%s %s/%s", BuildTime, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println("错误：", err)
	}
}
