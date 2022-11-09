package lockx

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type ReentrantLock struct {
	sync.Mutex
	recursion int32 // 这个goroutine 重入的次数
	owner     int64 // 当前持有锁的goroutine id
}

// GetGoroutineID
//
//	@Description: Get returns the id of the current goroutine.
//	@return int64
func GetGoroutineID() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return int64(n)
}

func NewReentrantLock() sync.Locker {
	res := &ReentrantLock{
		Mutex:     sync.Mutex{},
		recursion: 0,
		owner:     0,
	}
	return res
}

// ReentrantMutex 包装一个Mutex,实现可重入
type ReentrantMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的goroutine id
	recursion int32 // 这个goroutine 重入的次数
}

func (m *ReentrantMutex) Lock() {
	gid := GetGoroutineID()
	// 如果当前持有锁的goroutine就是这次调用的goroutine,说明是重入
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	// 获得锁的goroutine第一次调用，记录下它的goroutine id,调用次数加1
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *ReentrantMutex) Unlock() {
	gid := GetGoroutineID()
	// 非持有锁的goroutine尝试释放锁，错误的使用
	if atomic.LoadInt64(&m.owner) != gid {
		panic(fmt.Sprintf("wrong the owner(%d): %d!", m.owner, gid))
	}
	// 调用次数减1
	m.recursion--
	if m.recursion != 0 { // 如果这个goroutine还没有完全释放，则直接返回
		return
	}
	// 此goroutine最后一次调用，需要释放锁
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}
