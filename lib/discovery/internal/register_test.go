package internal

import (
	"git.zc0901.com/go/god/lib/logx"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var mockLock sync.Mutex

func init() {
	logx.Disable()
}

func setMockClient(cli EtcdClient) func() {
	mockLock.Lock()
	NewClient = func([]string) (EtcdClient, error) {
		return cli, nil
	}
	return func() {
		NewClient = DialClient
		mockLock.Unlock()
	}
}

func TestGetCluster(t *testing.T) {
	c1 := GetRegistry().getCluster([]string{"first"})
	c2 := GetRegistry().getCluster([]string{"second"})
	c3 := GetRegistry().getCluster([]string{"first"})
	assert.Equal(t, c1, c3)
	assert.NotEqual(t, c1, c2)
}

func TestGetClusterKey(t *testing.T) {
	assert.Equal(t,
		getClusterKey([]string{"localhost:1234", "remote-host:5678"}),
		getClusterKey([]string{"remote-host:5678", "localhost:1234"}))
}

func TestCluster_HandleChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	l := NewMockUpdateListener(ctrl)
	l.EXPECT().OnAdd(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnAdd(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnDelete(KV{
		Key: "first",
		Val: "1",
	})
	l.EXPECT().OnDelete(KV{
		Key: "second",
		Val: "2",
	})
	l.EXPECT().OnAdd(KV{
		Key: "third",
		Val: "3",
	})
	l.EXPECT().OnAdd(KV{
		Key: "fourth",
		Val: "4",
	})

	c := newCluster([]string{"any"})
	c.listeners["any"] = []UpdateListener{l}

	c.handleChanges("any", []KV{
		{
			"first",
			"1",
		},
		{
			"second",
			"2",
		},
	})
	assert.EqualValues(t, map[string]string{
		"first":  "1",
		"second": "2",
	}, c.values["any"])

	c.handleChanges("any", []KV{
		{
			"third",
			"3",
		},
		{
			"fourth",
			"4",
		},
	})
	assert.EqualValues(t, map[string]string{
		"third":  "3",
		"fourth": "4",
	}, c.values["any"])
}
