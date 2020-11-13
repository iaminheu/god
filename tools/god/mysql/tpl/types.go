package tpl

var Types = `
type (
	{{.table}} struct {
		{{.fields}}
	}

	{{.table}}Model struct {
		{{if .withCache}}sqlx.CachedConn{{else}}conn sqlx.Conn{{end}}
		table string
	}
)
`
