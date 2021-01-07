package tpl

var Update = `
func (m *{{.upperTable}}Model) Update(data {{.upperTable}}) error {
	{{if .withCache}}{{.primaryCacheKey}}
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + ` m.table +` + "` " + `set ` + "` + " + `{{.lowerTable}}FieldsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return conn.Exec(query, {{.values}})
	}, {{.primaryKeyName}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerTable}}FieldsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
	_,err := m.conn.Exec(query, {{.values}}){{end}}
	return err
}
`

var TxUpdate = `
func (m *{{.upperTable}}Model) TxUpdate(tx sqlx.TxSession, data {{.upperTable}}) error {
	{{if .withCache}}{{.primaryCacheKey}}
	_, err := m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + ` m.table +` + "` " + `set ` + "` + " + `{{.lowerTable}}FieldsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return tx.Exec(query, {{.values}})
	}, {{.primaryKeyName}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `{{.lowerTable}}FieldsWithPlaceHolder` + " + `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
	_,err := tx.Exec(query, {{.values}}){{end}}
	return err
}
`
