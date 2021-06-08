package testmodel

import (
	"fmt"
	"git.zc0901.com/go/god/lib/store/cache"
	"git.zc0901.com/go/god/lib/store/redis"
	"git.zc0901.com/go/god/lib/store/sqlx"
	"testing"
)

type Config struct {
	DataSource string
	Table      string
	Cache      cache.ClusterConf
}

type Model struct {
	c            Config
	ServiceModel *ServiceModel
}

func NewModel() *Model {
	c := Config{
		DataSource: "root:FfRyn2b5BKM3MNPz@tcp(dev:33061)/dci2?parseTime=true&charset=utf8mb4",
		Cache: cache.ClusterConf{
			{
				Conf: redis.Conf{
					// Host: "106.54.101.160:6382",
					Host: "192.168.0.17:6382",
					Mode: redis.StandaloneMode,
				},
				Weight: 100,
			},
		},
	}

	return &Model{
		c:            c,
		ServiceModel: NewServiceModel(sqlx.NewMySQL(c.DataSource), c.Cache),
	}
}

func TestServiceModel_SqlxNullXxx(t *testing.T) {
	m := NewModel()

	//s, err := m.ServiceModel.FindOne(1)
	//if err != nil {
	//	t.Error(err)
	//}
	//t.Log(s)

	services := m.ServiceModel.FindMany([]int64{1})
	for _, s := range services {
		fmt.Println(s.SettledAt)
	}
}
