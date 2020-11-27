package generator

import (
	"git.zc0901.com/go/god/lib/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	conf "git.zc0901.com/go/god/tools/god/config"
	"git.zc0901.com/go/god/tools/god/rpc/parser"
	"git.zc0901.com/go/god/tools/god/util"
	"git.zc0901.com/go/god/tools/god/util/format"
)

const configTemplate = `package config

import "git.zc0901.com/go/god/rpc"

type Config struct {
	rpc.ServerConf
}
`

func (g *defaultGenerator) GenConfig(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetConfig()
	configFilename, err := format.FileNamingFormat(cfg.NamingFormat, "config")
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, configFilename+".go")
	if fs.FileExist(fileName) {
		return nil
	}

	text, err := util.LoadTemplate(category, configTemplateFileFile, configTemplate)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(text), os.ModePerm)
}
