package gen

import (
	"git.zc0901.com/go/god/tools/god/mysql/tpl"
	"git.zc0901.com/go/god/tools/god/util"
)

func genTag(fieldName string) (string, error) {
	if fieldName == "" {
		return fieldName, nil
	}

	output, err := util.With("tag").Parse(tpl.Tag).Execute(map[string]interface{}{
		"field": fieldName,
	})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
