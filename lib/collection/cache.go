package collection

import (
	"container/list"
	"git.zc0901.com/go/god/lib/logx"
	"git.zc0901.com/go/god/lib/mathx"
	"git.zc0901.com/go/god/lib/syncx"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultCacheName = "proc" // 进程内缓存名
	slots            = 300
	statInterval     = time.Minute

	// 缓存过期偏差值：通过过期偏差，避免大量缓存同时过期
	// 将缓存到期时间偏差上下设置为0.05，也就是落到[0.95,1.05]*秒
	expireDeviation = 0.05
)

var emptyLruCache = emptyLru{}

type (
	Cache struct {
		name           string
		lock           sync.Mutex
		data           map[string]interface{}
		expire         time.Duration
		timingWheel    *TimingWheel
		lruCache       lru
		sharedCalls    syncx.SharedCalls
		unstableExpiry mathx.Unstable
		stats          *cacheStat
	}

	CacheOption func(cache *Cache)
)

// NewCache 新建一个进程内缓存
func NewCache(expire time.Duration, opts ...CacheOption) (*Cache, error) {
	cache := &Cache{
		data:           make(map[string]interface{}),
		expire:         expire,
		lruCache:       emptyLruCache,
		sharedCalls:    syncx.NewSharedCalls(),
		unstableExpiry: mathx.NewUnstable(expireDeviation),
	}

	for _, opt := range opts {
		opt(cache)
	}

	if len(cache.name) == 0 {
		cache.name = defaultCacheName
	}
	cache.stats = newCacheStat(cache.name, cache.size)

	timingWheel, err := NewTimingWheel(time.Second, slots, func(k, v interface{}) {
		key, ok := k.(string)
		if !ok {
			return
		}

		cache.Del(key)
	})
	if err != nil {
		return nil, err
	}

	cache.timingWheel = timingWheel
	return cache, nil
}

func WithLimit(limit int) CacheOption {
	return func(cache *Cache) {
		if limit > 0 {
			cache.lruCache = newKeyLru(limit, cache.onEvict)
		}
	}
}

func WithName(name string) CacheOption {
	return func(cache *Cache) {
		cache.name = name
	}
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	delete(c.data, key)
	c.lruCache.remove(key)
	c.lock.Unlock()
	c.timingWheel.RemoveTimer(key)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	value, ok := c.doGet(key)
	if ok {
		c.stats.IncrHit()
	} else {
		c.stats.IncrMiss()
	}

	return value, ok
}

func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	_, ok := c.data[key]
	c.data[key] = value
	c.lruCache.add(key)
	c.lock.Unlock()

	expiry := c.unstableExpiry.AroundDuration(c.expire)
	if ok {
		c.timingWheel.MoveTimer(key, expiry)
	} else {
		c.timingWheel.SetTimer(key, value, expiry)
	}
}

// Take 缓存已有则返回，反之则获取（通过共享调用实现高并发）
func (c *Cache) Take(key string, fetch func() (interface{}, error)) (interface{}, error) {
	if val, ok := c.doGet(key); ok {
		c.stats.IncrHit()
		return val, nil
	}

	// 回源获取，通过共享调用实现高并发
	val, hit, err := c.sharedCalls.Do(key, func() (interface{}, error) {
		// 因为在内存中进行的map搜索，时间为O(1)，而fetch是在IO上的查询
		// 所以我们做双重检测，因为缓存有可能被其他调用取到了
		if val, ok := c.doGet(key); ok {
			return val, nil
		}

		v, e := fetch()
		if e != nil {
			return nil, e
		}

		c.Set(key, v)
		return v, nil
	})
	if err != nil {
		return nil, err
	}

	if hit {
		c.stats.IncrHit()
	} else {
		c.stats.IncrMiss()
	}

	return val, nil
}

func (c *Cache) doGet(key string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.data[key]
	if ok {
		c.lruCache.add(key)
	}

	return value, ok
}

func (c *Cache) onEvict(key string) {
	delete(c.data, key)
	c.timingWheel.RemoveTimer(key)
}

func (c *Cache) size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

//// LRU

type (
	lru interface {
		add(key string)
		remove(key string)
	}

	emptyLru struct{}

	keyLru struct {
		limit    int
		evicts   *list.List
		elements map[string]*list.Element
		onEvict  func(key string)
	}
)

func newKeyLru(limit int, onEvict func(key string)) *keyLru {
	return &keyLru{
		limit:    limit,
		evicts:   list.New(),
		elements: make(map[string]*list.Element),
		onEvict:  onEvict,
	}
}

func (el emptyLru) add(string) {
}

func (el emptyLru) remove(string) {
}

func (kl *keyLru) add(key string) {
	if elem, ok := kl.elements[key]; ok {
		kl.evicts.MoveToFront(elem)
		return
	}

	// 添加新项
	elem := kl.evicts.PushFront(key)
	kl.elements[key] = elem

	// 防止下标溢出
	if kl.evicts.Len() > kl.limit {
		kl.removeOldest()
	}
}

func (kl *keyLru) remove(key string) {
	if elem, ok := kl.elements[key]; ok {
		kl.removeElement(elem)
	}
}

func (kl *keyLru) removeOldest() {
	elem := kl.evicts.Back()
	if elem != nil {
		kl.removeElement(elem)
	}
}

func (kl *keyLru) removeElement(e *list.Element) {
	kl.evicts.Remove(e)
	key := e.Value.(string)
	delete(kl.elements, key)
	kl.onEvict(key)
}

type cacheStat struct {
	name         string
	hit          uint64
	miss         uint64
	sizeCallback func() int
}

func newCacheStat(name string, sizeCallback func() int) *cacheStat {
	cs := &cacheStat{
		name:         name,
		sizeCallback: sizeCallback,
	}
	go cs.loop()
	return cs
}

func (cs *cacheStat) loop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for range ticker.C {
		hit := atomic.SwapUint64(&cs.hit, 0)
		miss := atomic.SwapUint64(&cs.miss, 0)
		total := hit + miss
		if total == 0 {
			continue
		}
		hitRatio := 100 * float32(hit) / float32(total)
		logx.Statf("cache(%s) - qpm: %d, hitRatio: %.1f%%, elements: %d, hit: %d, miss: %d",
			cs.name, total, hitRatio, cs.sizeCallback(), hit, miss)
	}
}

func (cs *cacheStat) IncrHit() {
	atomic.AddUint64(&cs.hit, 1)
}

func (cs *cacheStat) IncrMiss() {
	atomic.AddUint64(&cs.miss, 1)
}
