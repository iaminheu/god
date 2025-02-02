package clickhouse

import "git.zc0901.com/go/god/lib/store/sqlx"

// New 创建 ClickHouse 数据库实例
func New(dataSourceName string, opts ...sqlx.Option) sqlx.Conn {
	return sqlx.NewConn("clickhouse", dataSourceName, opts...)
}
