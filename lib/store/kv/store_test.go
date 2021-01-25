package kv

import (
	"git.zc0901.com/go/god/lib/hash"
	"git.zc0901.com/go/god/lib/store/cache"
	"git.zc0901.com/go/god/lib/store/redis"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	s1, _ = miniredis.Run()
	s2, _ = miniredis.Run()
)

func runOnCluster(t *testing.T, fn func(cluster Store)) {
	s1.FlushAll()
	s2.FlushAll()

	store := NewStore([]cache.Conf{
		{
			Conf: redis.Conf{
				Host: s1.Addr(),
				Mode: redis.StandaloneMode,
			},
			Weight: 100,
		},
		{
			Conf: redis.Conf{
				Host: s2.Addr(),
				Mode: redis.StandaloneMode,
			},
			Weight: 100,
		},
	})

	fn(store)
}

func TestRedis_Exists(t *testing.T) {
	store := clusterStore{dispatcher: hash.NewConsistentHash()}
	_, err := store.Exists("foo")
	assert.NotNil(t, err)

	runOnCluster(t, func(client Store) {
		ok, err := client.Exists("a")
		assert.Nil(t, err)
		assert.False(t, ok)
		assert.Nil(t, client.Set("a", "b"))
		ok, err = client.Exists("a")
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestRedis_Eval(t *testing.T) {
	store := clusterStore{dispatcher: hash.NewConsistentHash()}
	_, err := store.Eval(`redis.call("EXISTS", KEYS[1])`, "key1")
	assert.NotNil(t, err)

	runOnCluster(t, func(client Store) {
		_, err := client.Eval(`redis.call("EXISTS", KEYS[1])`, "notexist")
		assert.Equal(t, redis.Nil, err)
		err = client.Set("key1", "value1")
		assert.Nil(t, err)
		_, err = client.Eval(`redis.call("EXISTS", KEYS[1])`, "key1")
		assert.Equal(t, redis.Nil, err)
		val, err := client.Eval(`return redis.call("EXISTS", KEYS[1])`, "key1")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
	})
}
