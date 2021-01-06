package gen

import (
	"git.zc0901.com/go/god/lib/stringx"
	"git.zc0901.com/go/god/tools/god/mysql/tpl"
	"git.zc0901.com/go/god/tools/god/util"
)

func genFindMany(table Table) (string, error) {
	upperTable := table.Name.ToCamel()
	output, err := util.With("findMany").Parse(tpl.FindMany).Execute(map[string]interface{}{
		"upperTable": upperTable,
		"primaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).UnTitle(),
		"dataType":   table.PrimaryKey.DataType,
	})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
