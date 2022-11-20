package lru

type LruNode[K comparable, V any] struct {
	Key        K
	Val        V
	Prev, Next *LruNode[K, V]
}

type LruCache[K comparable, V any] struct {
	Size       int //节点的数量
	Cap        int // 容量
	Cache      map[K]*LruNode[K, V]
	Head, Tail *LruNode[K, V] //头尾节点
}

// NewLruCache[K comparable, V any]
//
//	@Description: 创建缓存对象
//	@param cap 缓存的数量
//	@return *LruCache[K，V]
func NewLruCache[K comparable, V any](cap int) *LruCache[K, V] {
	lruCache := &LruCache[K, V]{
		Size:  0,
		Cap:   cap,
		Cache: map[K]*LruNode[K, V]{},
		Head:  &LruNode[K, V]{},
		Tail:  &LruNode[K, V]{},
	}
	lruCache.Head.Next = lruCache.Tail
	lruCache.Tail.Prev = lruCache.Head
	return lruCache
}

// Get
//
//	@Description: 缓存中获取
//	@receiver l
//	@param key
//	@return V 命中返回值 未命中返回false和V类型的的零值
func (l *LruCache[K, V]) Get(key K) (bool, V) {
	node, ok := l.Cache[key]
	if !ok { //  Key不存在
		var v V
		return false, v
	}
	// 如果 key 存在，先通过哈希表定位，再移到头部
	l.moveToHead(node)
	return true, node.Val
}

// Put
//
//	@Description: 插入缓存
//	@receiver l
//	@param key 缓存的key
//	@param val 缓存的value
func (l *LruCache[K, V]) Put(key K, val V) {
	node, ok := l.Cache[key]
	if !ok { // 如果 key 不存在，创建一个新的节点
		newNode := &LruNode[K, V]{Key: key, Val: val}
		l.Cache[key] = newNode // 添加进哈希表
		l.addToHead(newNode)
		l.Size++
		if l.Size > l.Cap {
			// 如果超出容量，删除双向链表的尾部节点
			tail := l.removeTail()
			// 删除哈希表中对应的项
			delete(l.Cache, tail.Key)
			l.Size--
		}
	} else {
		node.Val = val
		// 如果 key 存在，先通过哈希表定位，再修改 value，并移到头部
		l.moveToHead(node)
	}
}

// addToHead
//
//	@Description: 添加到头部
//	@receiver l
//	@param node 待添加的节点
func (l *LruCache[K, V]) addToHead(node *LruNode[K, V]) {
	node.Prev = l.Head
	node.Next = l.Head.Next
	l.Head.Next.Prev = node
	l.Head.Next = node
}

// removeNode
//
//	@Description: 删除节点
//	@receiver l
//	@param node 待删除的节点
func (l *LruCache[K, V]) removeNode(node *LruNode[K, V]) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
}

// moveToHead
//
//	@Description: 把节点移动到头部
//	@receiver l
//	@param node 待移动的节点
func (l *LruCache[K, V]) moveToHead(node *LruNode[K, V]) {
	l.removeNode(node)
	l.addToHead(node)
}

// removeTail
//
//	@Description: 删除尾节点，先找到尾节点再删除
//	@receiver l
//	@return *LruNode[K,V] 被删除的尾节点
func (l *LruCache[K, V]) removeTail() *LruNode[K, V] {
	res := l.Tail.Prev
	l.removeNode(res)
	return res
}
