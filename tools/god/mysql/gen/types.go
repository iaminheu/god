package gen

import (
	"git.zc0901.com/go/god/tools/god/mysql/tpl"
	"git.zc0901.com/go/god/tools/god/util"
)

func genTypes(table Table, withCache bool) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields)
	if err != nil {
		return "", err
	}
	output, err := util.With("types").Parse(tpl.Types).Execute(map[string]interface{}{
		"withCache": withCache,
		"table":     table.Name.ToCamel(),
		"fields":    fieldsString,
	})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
