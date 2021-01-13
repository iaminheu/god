package builder

import (
	"fmt"
	"testing"
	"time"
)

func TestFieldList(t *testing.T) {
	type Dict struct {
		Id             int64     `db:"id"`        // 字典表 | 公共库
		ParentId       int64     `db:"parent_id"` // 字典分类ID
		Name           string    `db:"name"`      // 字典名称
		Type           string    `db:"type"`      // 字典类型
		CreateTime     time.Time // 创建时间
		CreateUserId   int64     `db:"create_user_id"`   // 创建人id
		CreateUserName string    `db:"create_user_name"` // 创建人姓名
		UpdateTime     time.Time `db:"update_time"`      // 更新时间
		UpdateUserId   int64     `db:"update_user_id"`   // 更新人id
		UpdateUserName string    `db:"update_user_name"` // 更新者姓名
		DeleteFlag     int64     `db:"delete_flag"`      // 删除标记: 0删除|1未删除
	}

	dict := &Dict{}
	resp := FieldList(dict)
	fmt.Println(resp)
}
