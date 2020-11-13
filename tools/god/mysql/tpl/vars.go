package tpl

var Vars = `
var (
	{{.table}}FieldList          = builder.FieldList(&{{.camelTable}}{})
	{{.table}}Fields                = strings.Join({{.table}}FieldList, ",")
	{{.table}}FieldsAutoSet         = strings.Join(stringx.Remove({{.table}}FieldList, {{if .autoIncrement}}"{{.primaryKey}}",{{end}} "created_at", "updated_at", "create_time", "update_time"), ",")
	{{.table}}FieldsWithPlaceHolder = strings.Join(stringx.Remove({{.table}}FieldList, "{{.primaryKey}}", "created_at", "updated_at", "create_time", "update_time"), "=?,") + "=?"

	{{if .withCache}}{{.cacheKeys}}{{end}}
)
`
