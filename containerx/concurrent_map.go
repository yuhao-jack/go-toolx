package containerx

import (
	"runtime"
	"sync"
	"unsafe"
)

type ConcurrentMapShared[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

var ShardCount = runtime.NumCPU() * 8

// ConcurrentMap 分成ShardCount个分片的map
type ConcurrentMap[K comparable, V any] []*ConcurrentMapShared[K, V]

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	m := make(ConcurrentMap[K, V], ShardCount)
	for i := 0; i < ShardCount; i++ {
		m[i] = &ConcurrentMapShared[K, V]{items: map[K]V{}}
	}
	return &m
}

// GetShard
//
//	@Description:
//	@receiver c
//	@param key
//	@return *ConcurrentMapShared[K
//	@return V]
func (c *ConcurrentMap[K, V]) GetShard(key K) *ConcurrentMapShared[K, V] {
	u := c.hash(key) % uintptr(ShardCount)
	return (*c)[int(u)]
}

// Set
//
//	@Description: 设置值，如果key存在，返回旧的值 否则返回V的零值
//	@receiver c
//	@param key
//	@param val
func (c *ConcurrentMap[K, V]) Set(key K, val V) (oldValue V) {
	shard := c.GetShard(key)
	shard.Lock()
	oldValue = shard.items[key]
	shard.items[key] = val
	shard.Unlock()
	return oldValue
}

// Get
//
//	@Description:
//	@receiver c
//	@param key
//	@param defaultVal
//	@return V
func (c *ConcurrentMap[K, V]) Get(key K, defaultVal ...V) (V, bool) {
	shard := c.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.items[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0], ok
		}
		var v V
		return v, ok
	}
	return val, ok
}

// Each
//
//	@Description:
//	@receiver c
//	@param f
func (c *ConcurrentMap[K, V]) Each(f func(key K, val V)) {
	for _, c2 := range *c {
		c2.Lock()
		for k, v := range c2.items {
			f(k, v)
		}
		c2.Unlock()
	}
}

// Remove
//
//	@Description:
//	@receiver c
//	@param key
func (c *ConcurrentMap[K, V]) Remove(key K) {
	shard := c.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// hash
//
//	@Description: 这里的hash方法是从一个帖子中找到的，测试过还不错
//	@See https://blog.csdn.net/weixin_45583158/article/details/106894015
//	@receiver c
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-05 17:25:42
//	@param key
//	@return uintptr
func (c *ConcurrentMap[K, V]) hash(key K) uintptr {
	var m interface{} = make(map[K]struct{})
	hf := (*mh)(*(*unsafe.Pointer)(unsafe.Pointer(&m))).hf
	return hf(unsafe.Pointer(&key), 0)
}

// mh is an inlined combination of runtime._type and runtime.maptype
type mh struct {
	_  uintptr
	_  uintptr
	_  uint32
	_  uint8
	_  uint8
	_  uint8
	_  uint8
	_  func(unsafe.Pointer, unsafe.Pointer) bool
	_  *byte
	_  int32
	_  int32
	_  unsafe.Pointer
	_  unsafe.Pointer
	_  unsafe.Pointer
	hf func(unsafe.Pointer, uintptr) uintptr
}
