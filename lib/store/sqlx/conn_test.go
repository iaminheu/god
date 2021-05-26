package sqlx

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Dict struct {
	Id             int64     `db:"id"`               // 字典表 | 公共库
	ParentId       int64     `db:"parent_id"`        // 字典分类ID
	Name           string    `db:"name"`             // 字典名称
	Type           string    `db:"type"`             // 字典类型
	CreateTime     time.Time `db:"create_time"`      // 创建时间
	CreateUserId   int64     `db:"create_user_id"`   // 创建人id
	CreateUserName string    `db:"create_user_name"` // 创建人姓名
	UpdateTime     time.Time `db:"update_time"`      // 更新时间
	UpdateUserId   int64     `db:"update_user_id"`   // 更新人id
	UpdateUserName string    `db:"update_user_name"` // 更新者姓名
	DeleteFlag     int64     `db:"delete_flag"`      // 删除标记: 0删除|1未删除
}

func TestABC(t *testing.T) {
	dataSourceName := "root:FfRyn2b5BKM3MNPz@tcp(192.168.0.17:33061)/nest_public?parseTime=true"
	db := NewMySQL(dataSourceName)

	var dictList []*Dict
	err := db.Query(&dictList, "select id, name, create_time from dict limit 0, 5")
	if err != nil {
		panic(err)
	}

	for _, dict := range dictList {
		fmt.Println(dict.Id, dict.Name, dict.CreateTime)
	}
}

func TestDbInstance_QueryRows(t *testing.T) {
	dataSourceName := "root:FfRyn2b5BKM3MNPz@tcp(192.168.0.17:33061)/nest_label?parseTime=true"
	db := NewMySQL(dataSourceName)
	type AccountKinds []struct {
		Id   int
		Name string
	}

	var book struct {
		Name  string
		Total int
		Price float32
		kinds AccountKinds
	}
	type Books []struct {
		Total         int    `conn:"totalx"`
		Name          string `conn:"book"`
		NotExistField int    `conn:"y"`
	}

	var accountKinds AccountKinds
	var books Books
	var adminUsers []struct {
		Txt       string    `db:"txt"`
		UserId    int       `db:"user_id"`
		AdminId   int       `db:"admin_id"`
		CreatedAt time.Time `db:"created_at"`
	}

	// 查询测试
	errAccountKinds := db.Query(&accountKinds, "select id, value as name from nest_user.account_kind")
	// errBook := db.Query(&book, "select book, count(0) total from book group by book order by total desc")
	// errBooks := db.Query(&books, "select book, count(0) totalx, 1 as x, 2 as y from book group by book order by totalx desc")
	errAdminUsers := db.Query(&adminUsers, "select user_id, admin_id, txt, created_at from nest_admin.admin_user")

	if errAccountKinds != nil {
		t.Fatal(errAccountKinds)
	}

	//if errBook != nil {
	//	t.Fatal(errBook)
	//}
	//
	//if errBooks != nil {
	//	t.Fatal(errBooks)
	//}

	book.kinds = accountKinds

	if errAdminUsers != nil {
		t.Fatal(errAdminUsers)
	}

	fmt.Println(book)

	for _, book := range books {
		fmt.Println(book)
	}

	for _, accountKind := range accountKinds {
		fmt.Println(accountKind)
	}

	for _, adminUser := range adminUsers {
		fmt.Println(adminUser)
	}
}

func TestDbInstance_Exec(t *testing.T) {
	dataSourceName := "root:FfRyn2b5BKM3MNPz@tcp(192.168.0.17:33061)/nest_label?parseTime=true"
	db := NewMySQL(dataSourceName)

	for i := 1; i <= 7; i++ {
		result, err := db.Exec("update nest_admin.admin_user set txt=? where id=?", "自在测试"+strconv.Itoa(i), i)
		if err != nil {
			t.Fatal(err)
		}
		lastInsertId, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("LastInsertId: %d, RowsAffected: %d\n", lastInsertId, rowsAffected)
	}
}

func TestConn_Prepare(t *testing.T) {
	dataSourceName := "root:FfRyn2b5BKM3MNPz@tcp(192.168.0.17:33061)/nest_label?parseTime=true"
	db := NewMySQL(dataSourceName)

	query := "update nest_admin.admin_user set txt=? where id=?"
	stmt, err := db.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	defer stmt.Close()

	for i := 1; i <= 7; i++ {
		result, err := stmt.Exec("测试"+strconv.Itoa(i), i)
		if err != nil {
			t.Fatal(err)
		}
		lastInsertId, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("LastInsertId: %d, RowsAffected: %d\n", lastInsertId, rowsAffected)
	}
}

func TestConnBreaker(t *testing.T) {
	// logx.Disable()
	// logx.SetLevel(logx.ErrorLevel)
	// dataSourceName := "root:asdfasdf@tcp(192.168.0.166:33061)/nest_label?parseTime=true&timeout=10s&readTimeout=2s"
	// dataSourceName := "root:asdfasdf@tcp(218.244.143.31:3317)/nest_label?parseTime=true&timeout=1s&readTimeout=2s"
	dataSourceName := "root:asdfasdf@tcp(218.244.143.31:3317)/nest_label?parseTime=true"
	db := NewMySQL(dataSourceName)
	var book struct {
		Book string `conn:"book"`
	}

	for i := 0; i < 100; i++ {
		_ = db.Query(&book, "select book from bookx limit ?", i)
	}
}

func Test_Scan(t *testing.T) {
	sqlGetMenuByRoleId := `select m.id, m.parent_id, m.name, m.title, m.path, m.component, m.icon, m.keep_alive, hidden, update_time
	from menu m
	right join role_menu rm on m.id=rm.menu_id
	where rm.role_id=? order by m.sort`

	type MenuResp struct {
		Title      string       `db:"title"`             // 菜单名
		Name       string       `db:"name"`              // 菜单名
		UpdateTime sql.NullTime `db:"update_time"`       // 更新时间
		Children   []*MenuResp  `db:"-" json:"children"` // （不属于sql的字段，一定要声明db:"-"，进行忽略）
	}

	var roleMenus []*MenuResp
	dataSourceName := "root:FfRyn2b5BKM3MNPz@tcp(106.54.101.160:33061)/nest_casbin?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	db := NewMySQL(dataSourceName)

	err := db.Query(&roleMenus, sqlGetMenuByRoleId, 1)
	if err == nil {
		for _, menu := range roleMenus {
			fmt.Println(menu.Title, menu.Name, menu.UpdateTime)
		}
	}
}
