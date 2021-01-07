package tpl

var Insert = `
func (m *{{.upperTable}}Model) Insert(data {{.upperTable}}) (sql.Result, error) {
	query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerTable}}FieldsAutoSet` + " + `) values ({{.args}})` " + `
	return m.{{if .withCache}}ExecNoCache{{else}}conn.Exec{{end}}(query, {{.values}})
}
`

var TxInsert = `
func (m *{{.upperTable}}Model) TxInsert(tx sqlx.TxSession, data {{.upperTable}}) (sql.Result, error) {
	query := ` + "`" + `insert into ` + "`" + ` + m.table + ` + "` (` + " + `{{.lowerTable}}FieldsAutoSet` + " + `) values ({{.args}})` " + `
	return tx.Exec(query, {{.values}})
}
`
