package lfu

import (
	"sort"
	"time"
)

type K interface {
	comparable
}

type LfuCache[K comparable, V any] struct {
	capcity int
	cache   map[K]V
	count   map[K]*HitRate[K]
}

type HitRate[K comparable] struct {
	key      K
	hitCount int
	lastTime int64
}

// NewLruCache[K comparable, V any]
//
//	@Description: 创建LFU缓存对象
//	@param capcity 缓存容量
//	@return *LfuCache[K, V]
func NewLfuCache[K comparable, V any](capcity int) *LfuCache[K, V] {
	lfuCache := &LfuCache[K, V]{
		capcity: capcity,
		cache:   map[K]V{},
		count:   map[K]*HitRate[K]{},
	}
	return lfuCache
}

// Put
//
//	@Description: 插入缓存
//	@receiver l
//	@param key 缓存的key
//	@param val 缓存的val
func (l *LfuCache[K, V]) Put(key K, val V) {
	_, ok := l.cache[key]
	if !ok {
		if len(l.cache) == l.capcity {
			l.removeElement()
		}
		l.count[key] = &HitRate[K]{key: key, hitCount: 1, lastTime: time.Now().Unix()}
	} else {
		l.addHitCount(key)
	}
	l.cache[key] = val
}

// Get
//
//	@Description: 查询缓存
//	@receiver l
//	@param key 缓存key
//	@return V 缓存的val 不存在时返回V的零值
func (l *LfuCache[K, V]) Get(key K, defaultVal ...V) V {
	v, ok := l.cache[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		var v V
		return v
	}
	l.addHitCount(key)
	return v
}

// removeElement
//
//	@Description: 删除元素
//	@receiver l
func (l *LfuCache[K, V]) removeElement() {
	var hArr []*HitRate[K]
	for _, h := range l.count {
		hArr = append(hArr, h)
	}
	//  先按照命中次数排序
	//  命中次数相同再按照时间排序
	sort.Slice(hArr, func(i, j int) bool {
		if hArr[i].hitCount < hArr[j].hitCount {
			return true
		}
		if hArr[i].hitCount > hArr[j].hitCount {
			return false
		}
		return hArr[i].lastTime < hArr[j].lastTime
	})
	if len(hArr) > 0 {
		delete(l.cache, hArr[0].key)
		delete(l.count, hArr[0].key)
	}

}

// addHitCount
//
//	@Description: 更新访问元素状态
//	@receiver l
//	@param key
func (l *LfuCache[K, V]) addHitCount(key K) {
	h := l.count[key]
	h.hitCount++
	h.lastTime = time.Now().Unix()
}
