package generator

import (
	"git.zc0901.com/go/god/lib/fs"
	"path/filepath"
	"strings"

	"git.zc0901.com/go/god/lib/stringx"
	"git.zc0901.com/go/god/tools/god/util"
)

const rpcTemplateText = `syntax = "proto3";

package {{.package}};

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

service {{.serviceName}} {
  rpc Ping(Request) returns(Response);
}
`

func ProtoTmpl(out string) error {
	protoFilename := filepath.Base(out)
	serviceName := stringx.From(strings.TrimSuffix(protoFilename, filepath.Ext(protoFilename)))
	text, err := util.LoadTemplate(category, rpcTemplateFile, rpcTemplateText)
	if err != nil {
		return err
	}

	dir := filepath.Dir(out)
	err = fs.MkdirIfNotExist(dir)
	if err != nil {
		return err
	}

	err = util.With("t").Parse(text).SaveTo(map[string]string{
		"package":     serviceName.UnTitle(),
		"serviceName": serviceName.Title(),
	}, out, false)
	return err
}
