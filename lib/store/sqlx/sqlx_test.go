package sqlx

import (
	"fmt"
	"git.zc0901.com/go/god/lib/store/cache"
	"git.zc0901.com/go/god/lib/store/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Config struct {
	DataSource string
	Table      string
	Cache      cache.ClusterConf
}

type Model struct {
	c       Config
	Profile *ProfileModel
}

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

func NewModel() *Model {
	c := Config{
		DataSource: "root:qxqgqzx2018@tcp(106.54.101.160:3306)/nest_statistics?parseTime=true",
		Cache: cache.ClusterConf{
			{
				Conf: redis.Conf{
					Host: "106.54.101.160:6382",
					//Host: "192.168.0.17:6379",
					Mode: redis.StandaloneMode,
				},
				Weight: 100,
			},
		},
	}

	return &Model{
		c:       c,
		Profile: NewProfileModel(NewMySQL(c.DataSource)),
	}
}

func TestScan2Struct(t *testing.T) {
	model := NewModel()
	profile, err := model.Profile.FindOne(1)
	assert.Nil(t, err)
	fmt.Println(profile)
	fmt.Println(profile.Nickname)
}
