package fifo

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/algorithm/cache/fifo"
	"testing"
)

func TestFifo(t *testing.T) {
	cache := fifo.NewFifoCache[string, any](3)
	cache.Put("A", 23)
	fmt.Println(cache.String())
	cache.Put("B", 32)
	fmt.Println(cache.String())
	cache.Put("C", 223)
	fmt.Println(cache.String())
	cache.Put("D", 233)
	fmt.Println(cache.String())
	cache.Put("A", 233)
	fmt.Println(cache.String())

}
