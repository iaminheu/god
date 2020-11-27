package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"git.zc0901.com/go/god/lib/stringx"
	conf "git.zc0901.com/go/god/tools/god/config"
	"git.zc0901.com/go/god/tools/god/rpc/parser"
	"git.zc0901.com/go/god/tools/god/util"
	"git.zc0901.com/go/god/tools/god/util/format"
)

const etcTemplate = `Name: {{.serviceName}}.rpc
ListenOn: 127.0.0.1:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: {{.serviceName}}.rpc
`

func (g *defaultGenerator) GenEtc(ctx DirContext, _ parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetEtc()
	etcFilename, err := format.FileNamingFormat(cfg.NamingFormat, ctx.GetMain().Base)
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%v.yaml", etcFilename))

	text, err := util.LoadTemplate(category, etcTemplateFileFile, etcTemplate)
	if err != nil {
		return err
	}

	return util.With("etc").Parse(text).SaveTo(map[string]interface{}{
		"serviceName": strings.ToLower(stringx.From(ctx.GetMain().Base).ToCamel()),
	}, fileName, false)
}
