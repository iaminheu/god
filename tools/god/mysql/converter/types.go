package converter

import (
	"fmt"
	"strings"
)

var types = map[string]string{
	// 图省事，所有数据库整型字段都转为 int64
	// number
	"bool":      "int64",
	"boolean":   "int64",
	"tinyint":   "int64",
	"smallint":  "int64",
	"mediumint": "int64",
	"int":       "int64",
	"integer":   "int64",
	"bigint":    "int64",
	"float":     "float64",
	"double":    "float64",
	"decimal":   "float64",

	// date&time
	"date":      "time.Time",
	"datetime":  "time.Time",
	"timestamp": "time.Time",
	"time":      "string",
	"year":      "int64",

	// string
	"char":       "string",
	"varchar":    "string",
	"binary":     "string",
	"varbinary":  "string",
	"tinytext":   "string",
	"text":       "string",
	"mediumtext": "string",
	"longtext":   "string",
	"enum":       "string",
	"set":        "string",
	"json":       "string",
}

// ConvertDataType {
func ConvertDataType(dataBaseType string, isDefaultNull bool) (string, error) {
	tp, ok := types[strings.ToLower(dataBaseType)]
	if !ok {
		return "", fmt.Errorf("不识别的数据库类型: %s", dataBaseType)
	}
	return mayConvertNullType(tp, isDefaultNull), nil
}

func mayConvertNullType(goDataType string, isDefaultNull bool) string {
	if !isDefaultNull {
		return goDataType
	}

	switch goDataType {
	case "int64":
		return "sqlx.NullInt64"
	case "int32":
		return "sqlx.NullInt32"
	case "float64":
		return "sqlx.NullFloat64"
	case "bool":
		return "sqlx.NullBool"
	case "string":
		return "sqlx.NullString"
	case "time.Time":
		return "sqlx.NullTime"
	default:
		return goDataType
	}
}
