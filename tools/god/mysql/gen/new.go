package gen

import (
	"git.zc0901.com/go/god/tools/god/mysql/tpl"
	"git.zc0901.com/go/god/tools/god/util"
)

func genNew(table Table, database string, withCache bool) (string, error) {
	output, err := util.With("new").Parse(tpl.New).Execute(map[string]interface{}{
		"withCache":   withCache,
		"database":    database,
		"originTable": table.Name.Source(),
		"table":       table.Name.ToCamel(),
	})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
