package gen

import (
	"git.zc0901.com/go/god/tools/god/mysql/tpl"
	"git.zc0901.com/go/god/tools/god/util"
)

func genImports(withCache, timeImport bool) (string, error) {
	if withCache {
		buffer, err := util.With("import").Parse(tpl.Imports).Execute(map[string]interface{}{
			"time": timeImport,
		})
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	} else {
		buffer, err := util.With("import").Parse(tpl.ImportsNoCache).Execute(map[string]interface{}{
			"time": timeImport,
		})
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
}
