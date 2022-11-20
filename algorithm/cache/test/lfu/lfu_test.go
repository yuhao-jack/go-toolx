package lfu

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/algorithm/cache/lfu"
	"github.com/yuhao-jack/go-toolx/fun"
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
	ok, v := lfuCache.Get("A")
	fmt.Printf("%v\n", fun.IfOr(ok, v, -1))

}
