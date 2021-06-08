package sqlx

import (
	"fmt"
	"testing"
)

type Area struct {
	Id   int64  `conn:"id"`
	Code string `conn:"code"`
	Name string `conn:"name"`
}

func TestGetCityList(t *testing.T) {
	dataSourceName := "root:qxqgqzx2018@tcp(192.168.0.17:3306)/nest_public?parseTime=true"
	db := NewMySQL(dataSourceName)

	var areaList []*Area
	query := "select id, name, code from area where  parent_code = ?"
	err := db.Query(&areaList, query, 110100)
	if err != nil {
		fmt.Println(err)
	}
	for _, area := range areaList {
		fmt.Println(area.Id, area.Name, area.Code)
		fmt.Println(struct {
			ID   int
			City string
			No   string
		}{
			ID:   int(area.Id),
			City: area.Name,
			No:   area.Code,
		})
	}
}

func TestGreatThan(t *testing.T) {
	dataSourceName := "root:qxqgqzx2018@tcp(106.54.101.160:3306)/nest_statistics?parseTime=true"
	db := NewMySQL(dataSourceName)

	query := "update daily_account_num set total=10000 where total > 10000"
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
}

func TestSqlIn(t *testing.T) {
	ids := []int{2, 3}

	query := fmt.Sprintf("select id from user where id in (%s)", In(len(ids)))
	fmt.Println(query)
}

//func TestScan2Struct(t *testing.T) {
//	m := NewModel()
//
//	//type member struct {
//	//	Id   int64  `db:"id"`
//	//	Name string `db:"name"`
//	//}
//	//var members []*member
//	//
//	//err := m.ChannelContentModel.QueryNoCache(&members, "select id,name from ca_member limit 3")
//	//if err != nil {
//	//	t.Error(err)
//	//}
//	//
//	//for _, m := range members {
//	//	fmt.Printf("%d, %s\n", m.Id, m.Name)
//	//}
//
//	result, err := m.ChannelContentModel.Insert(ChannelContent{
//		ChannelId:   4,
//		ContentType: 4,
//		ContentId:   4,
//		ToppedAt:    NullTime{},
//		UnqKey:      "4",
//	})
//	if err != nil {
//		t.Error(err)
//	}
//
//	t.Log(result)
//}
