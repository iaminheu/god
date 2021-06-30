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

var UpdatePartial = `
func (m *{{.upperTable}}Model) UpdatePartial(data g.Map) error {
	updateArgs, err := sqlx.ExtractUpdateArgs({{.lowerTable}}FieldList, data)
	if err != nil {
		return err
	}

	{{if .withCache}}{{.primaryCacheKey}}
	_, err = m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := ` + "`" + `update ` + "` +" + ` m.table +` + "` " + `set ` + "` + " + `updateArgs.Fields` + " + `" + ` where {{.originalPrimaryKey}} = updateArgs.Id` + "`" + `
		return conn.Exec(query, updateArgs.Args...)
	}, {{.primaryKeyName}}){{else}}query := ` + "`" + `update ` + "` +" + `m.table +` + "` " + `set ` + "` +" + `updateArgs.Fields` + " + `" + ` where {{.originalPrimaryKey}} = updateArgs.Id` + "`" + `
	_,err := m.conn.Exec(query, updateArgs.Args...){{end}}
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
