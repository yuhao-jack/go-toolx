package lockx

import (
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

type ConcurrentMapShared[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

var ShardCount = runtime.NumCPU()

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
	var bytes []byte
	if reflect.TypeOf(key).Kind() != reflect.Ptr {
		bytes = *(*[]byte)(unsafe.Pointer(&key))
	} else {
		pointer := reflect.ValueOf(key).Pointer()
		bytes = *(*[]byte)(unsafe.Pointer(pointer))
	}
	return (*c)[fnv32(bytes)%ShardCount]
}

// Set
//
//	@Description:
//	@receiver c
//	@param key
//	@param val
func (c *ConcurrentMap[K, V]) Set(key K, val V) {
	shard := c.GetShard(key)
	shard.Lock()
	shard.items[key] = val
	shard.Unlock()
}

// Get
//
//	@Description:
//	@receiver c
//	@param key
//	@param defaultVal
//	@return V
func (c *ConcurrentMap[K, V]) Get(key K, defaultVal ...V) V {
	shard := c.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()
	val, ok := shard.items[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		var v V
		return v
	}
	return val
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

var Prime = 16777619
var OffsetBasis = 2166136261

func fnv32(src []byte) int {
	hash := OffsetBasis
	for _, b := range src {
		hash ^= int(b)
		hash *= Prime
	}
	return hash
}
