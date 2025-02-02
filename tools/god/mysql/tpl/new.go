package tpl

var New = `
func New{{.table}}Model(conn sqlx.Conn,{{if .withCache}} clusterConf cache.ClusterConf{{end}}) *{{.table}}Model {
	return &{{.table}}Model {
		{{if .withCache}}CachedConn: sqlx.NewCachedConnWithCluster(conn, clusterConf){{else}}conn: conn{{end}},
		table: "{{.originTable}}",
	}
}
`
