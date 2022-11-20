package lfu

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/algorithm/cache/lfu"
	"testing"
)

func TestLfu(t *testing.T) {
	lfuCache := lfu.NewLfuCache[string, int](3)
	lfuCache.Put("A", 1)
	lfuCache.Put("B", 2)
	lfuCache.Put("C", 3)
	lfuCache.Put("D", 2)
	lfuCache.Put("A", 1)
	lfuCache.Put("B", 5)
	v := lfuCache.Get("A", -1)
	fmt.Printf("%v\n", v)
	v = lfuCache.Get("F", -1)
	fmt.Printf("%v\n", v)
}
