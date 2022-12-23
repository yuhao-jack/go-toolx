package test

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/containerx"
	"testing"
	"unsafe"
)

type A struct {
	Nmae string
	Age  int
}

func TestDemo1(t *testing.T) {

	concurrentMap := containerx.NewConcurrentMap[string, string]()
	concurrentMap.Set("A", "A")
	concurrentMap.Set("a", "a")

}

func BenchmarkDemo1(b *testing.B) {
	a := &A{
		Nmae: "99",
		Age:  100,
	}
	for i := 0; i < b.N; i++ {
		fmt.Println(hash(a))
	}

}

//func hash[K comparable](key K) uintptr {
//	m := make(map[K]struct{})
//	hf := (*mh)(*(*unsafe.Pointer)(unsafe.Pointer(&m))).hf
//	return hf(unsafe.Pointer(&key), 0)
//}

//type mh struct {
//	_  uintptr
//	_  uintptr
//	_  uint32
//	_  uint8
//	_  uint8
//	_  uint8
//	_  uint8
//	_  func(unsafe.Pointer, unsafe.Pointer) bool
//	_  *byte
//	_  int32
//	_  int32
//	_  unsafe.Pointer
//	_  unsafe.Pointer
//	_  unsafe.Pointer
//	hf func(unsafe.Pointer, uintptr) uintptr
//}

func hash[A comparable](a A) uintptr {
	var m interface{} = make(map[A]struct{})
	hf := (*mh)(*(*unsafe.Pointer)(unsafe.Pointer(&m))).hf
	return hf(unsafe.Pointer(&a), 0)
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
