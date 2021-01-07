package tpl

var Delete = `
func (m *{{.upperStartCamelObject}}Model) Delete({{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return conn.Exec(query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		_,err:=m.conn.Exec(query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
`

var TxDelete = `
func (m *{{.upperStartCamelObject}}Model) TxDelete(tx sqlx.TxSession, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne({{.lowerStartCamelPrimaryKey}})
	if err!=nil{
		return err
	}{{end}}

	{{.keys}}
    _, err {{if .containsIndexCache}}={{else}}:={{end}} m.Exec(func(conn sqlx.Conn) (result sql.Result, err error) {
		query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		return tx.Exec(query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}query := ` + "`" + `delete from ` + "` +" + ` m.table + ` + " `" + ` where {{.originalPrimaryKey}} = ?` + "`" + `
		_,err:=tx.Exec(query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}
`
