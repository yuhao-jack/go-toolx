package ring_buf

import "github.com/yuhao-jack/go-toolx/fun"

type RingBufCache[T any] struct {
	// 当head==tail时，说明buffer为空，当head==(tail+1)%bufferSize则说明buffer满了
	head       int //指向下一次读的位置
	tail       int //指向的是下一次写的位置
	buffer     []T //数组来存储
	bufferSize int
}

// NewRingBufCache [T any]
//
//	@Description: 创建一个NewRingBufCache
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 15:31:50
//	@param bufferSize 大小
//	@return *RingBufCache[T]
func NewRingBufCache[T any](bufferSize int) *RingBufCache[T] {
	return &RingBufCache[T]{
		head:       0,
		tail:       0,
		buffer:     make([]T, bufferSize),
		bufferSize: bufferSize,
	}
}

// IsEmpty
//
//	@Description: 是否为空
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 14:54:31
//	@return bool
func (r *RingBufCache[T]) IsEmpty() bool {
	return r.head == r.tail
}

// IsFull
//
//	@Description: 是否满了
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 14:58:56
//	@return bool
func (r *RingBufCache[T]) IsFull() bool {
	return (r.tail+1)%r.bufferSize == r.head
}

// Clean
//
//	@Description: 清空ring_buf
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 15:05:17
func (r *RingBufCache[T]) Clean() {
	var t T
	for i := 0; i < r.bufferSize; i++ {
		r.buffer[i] = t
	}
	r.tail = 0
	r.head = 0
}

// Put
//
//	@Description: 插入数据
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 15:09:07
//	@param val
//	@return bool 满了则插入失败
func (r *RingBufCache[T]) Put(val T) bool {
	if r.IsFull() {
		return false
	}
	r.buffer[r.tail] = val
	r.tail = (r.tail + 1) % r.bufferSize
	return true
}

// Get
//
//	@Description: 获取数据，可以传入默认值，当传了默认值buf为空时返回传入的第一个默认值，否则返回T的零值
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 15:13:46
//	@param defaultVal
//	@return T
//	@return bool 从buf中取出的值为true 否则为false
func (r *RingBufCache[T]) Get(defaultVal ...T) (T, bool) {
	if r.IsEmpty() {
		if len(defaultVal) > 0 {
			return defaultVal[0], false
		}
		var t T
		return t, false
	}
	t := r.buffer[r.head]
	r.head = (r.head + 1) % r.bufferSize
	return t, true
}

// GetAll
//
//	@Description: 获取有效的数据
//	@receiver r
//	@Author yuhao <yuhao@mini1.cn>
//	@Data 2022-12-23 15:22:36
//	@return []T
func (r *RingBufCache[T]) GetAll() []T {
	if r.IsEmpty() {
		return []T{}
	}
	copyTail := r.tail
	cnt := fun.IfOr(r.head < copyTail, copyTail-r.head, r.bufferSize-r.head+copyTail)
	result := make([]T, cnt)
	if r.head < copyTail {
		for i := r.head; i < copyTail; i++ {
			result[i-r.head] = r.buffer[i]
		}
	} else {
		for i := r.head; i < r.bufferSize; i++ {
			result[i-r.head] = r.buffer[i]
		}
		for i := 0; i < copyTail; i++ {
			result[r.bufferSize-r.head+i] = r.buffer[i]
		}
	}
	r.head = copyTail
	return result
}
