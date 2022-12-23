package fifo

import (
	"fmt"
	"strings"
)

type FifoCache[K comparable, V any] struct {
	capacity   int
	size       int
	cache      map[K]*FifoNode[K, V]
	head, tail *FifoNode[K, V]
}

type FifoNode[K comparable, V any] struct {
	Key        K
	Val        V
	Prev, Next *FifoNode[K, V]
}

// NewFifoCache [K comparable, V any]
//
//	@Description: 创建缓存对象
//	@param capacity 缓存的数量
//	@return *FifoCache[K, V]
func NewFifoCache[K comparable, V any](capacity int) *FifoCache[K, V] {
	fifoCache := FifoCache[K, V]{
		capacity: capacity,
		size:     0,
		cache:    map[K]*FifoNode[K, V]{},
		head:     &FifoNode[K, V]{},
		tail:     &FifoNode[K, V]{},
	}

	fifoCache.head.Next = fifoCache.tail
	fifoCache.tail.Prev = fifoCache.head
	return &fifoCache
}

// String
//
//	@Description:
//	@receiver l
//	@return string
func (l *FifoCache[K, V]) String() string {
	if l == nil {
		return ""
	}
	sb := strings.Builder{}
	for _, f := range l.cache {
		if sb.Len() == 0 {
			sb.WriteString(fmt.Sprint("{", f.Key, "=", f.Val))
		} else {
			sb.WriteString(fmt.Sprint(",", f.Key, "=", f.Val))
		}
	}
	sb.WriteString("}")
	return sb.String()
}

// Get
//
//	@Description: 缓存中获取
//	@receiver l
//	@param key
//	@param defaultVal 未命中返回V类型的的零值
//	@return V 命中返回值 未命中返回V类型的的零值
func (l *FifoCache[K, V]) Get(key K, defaultVal ...V) (V, bool) {
	node, ok := l.cache[key]
	if !ok { //  Key不存在
		if len(defaultVal) > 0 {
			return defaultVal[0], ok
		}
		var v V
		return v, ok
	}
	return node.Val, ok
}

// Put
//
//	@Description: 插入缓存
//	@receiver l
//	@param key 缓存的key
//	@param val 缓存的value
func (l *FifoCache[K, V]) Put(key K, val V) {
	node, ok := l.cache[key]
	if !ok { // 如果 key 不存在，创建一个新的节点
		newNode := &FifoNode[K, V]{Key: key, Val: val}
		l.cache[key] = newNode
		l.addToHead(newNode)
		l.size++
	} else {
		node.Val = val
		l.moveToHead(node)
	}
	if l.size > l.capacity {
		tail := l.removeTail()
		delete(l.cache, tail.Key)
	}
}

// addToHead
//
//	@Description: 添加到头部
//	@receiver l
//	@param node 待添加的节点
func (l *FifoCache[K, V]) addToHead(node *FifoNode[K, V]) {
	node.Prev = l.head
	node.Next = l.head.Next
	l.head.Next.Prev = node
	l.head.Next = node
}

// removeTail
//
//	@Description: 删除尾节点，先找到尾节点再删除
//	@receiver l
//	@return *LruNode[K,V] 被删除的尾节点
func (l *FifoCache[K, V]) removeTail() *FifoNode[K, V] {
	res := l.tail.Prev
	l.removeNode(res)
	return res
}

// removeNode
//
//	@Description: 删除节点
//	@receiver l
//	@param node 待删除的节点
func (l *FifoCache[K, V]) removeNode(node *FifoNode[K, V]) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

// moveToHead
//
//	@Description: 把节点移动到头部
//	@receiver l
//	@param node 待移动的节点
func (l *FifoCache[K, V]) moveToHead(node *FifoNode[K, V]) {
	l.removeNode(node)
	l.addToHead(node)
}
