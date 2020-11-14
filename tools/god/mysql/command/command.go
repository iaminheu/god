package command

import (
	"errors"
	"git.zc0901.com/go/god/lib/collection"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/store/sqlx"
	"git.zc0901.com/go/god/tools/god/mysql/gen"
	"git.zc0901.com/go/god/tools/god/mysql/model"
	"git.zc0901.com/go/god/tools/god/util"
	"github.com/urfave/cli"
	"strings"
)

const (
	flagDSN   = "dsn"
	flagTable = "table"
	flagDir   = "dir"
	flagCache = "cache"
)

func GenCodeFromDSN(ctx *cli.Context) error {
	dsn := ctx.String(flagDSN)
	dir := ctx.String(flagDir)
	cache := ctx.Bool(flagCache)
	table := strings.TrimSpace(ctx.String(flagTable))

	logx.Disable()
	log := util.NewConsole(true)

	if len(dsn) == 0 {
		log.Error("MySQL连接地址未提供")
		return nil
	}
	if len(table) == 0 {
		log.Error("表名未提供")
		return nil
	}

	tables := collection.NewSet()
	for _, table := range strings.Split(table, ",") {
		table = strings.TrimSpace(table)
		if len(table) == 0 {
			continue
		}
		tables.AddStr(table)
	}

	conn := sqlx.NewMySQL(dsn)
	m := model.NewModel(conn)
	ddlList, err := m.ShowDDL(tables.KeysStr()...)
	if err != nil {
		log.Error("", err)
		return nil
	}

	// 获取数据库名称
	path := strings.Split(dsn, "?")[0]
	parts := strings.Split(path, "/")
	database := strings.TrimSpace(parts[len(parts)-1])
	if !strings.Contains(path, "/") || database == "" {
		log.Error("数据库连接字符串：未提供数据库名称")
		return errors.New("数据库连接字符串：未提供数据库名称")
	}

	//fmt.Println(strings.Join(ddlList, "\n"), dir, cache)
	generator := gen.NewModelGenerator(ddlList, dir, gen.WithConsoleOption(log))
	err = generator.Start(database, cache)
	if err != nil {
		log.Error("", err)
	}

	return nil
}
