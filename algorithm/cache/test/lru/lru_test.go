package lru_test

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/algorithm/cache/lru"
	"testing"
)

type User struct {
	Name string
	Age  int
	Sex  string
}

func TestLru1(t *testing.T) {
	lruCache := lru.NewLruCache[string, *User](3)
	lruCache.Put("A", &User{
		Name: "张三A",
		Age:  10,
		Sex:  "男",
	})
	lruCache.Put("B", &User{
		Name: "张三A",
		Age:  11,
		Sex:  "男",
	})
	lruCache.Put("C", &User{
		Name: "张三A",
		Age:  12,
		Sex:  "男",
	})
	lruCache.Put("D", &User{
		Name: "张三A",
		Age:  13,
		Sex:  "男",
	})
	lruCache.Put("A", &User{
		Name: "张三A",
		Age:  14,
		Sex:  "男",
	})
	lruCache.Put("B", &User{
		Name: "张三A",
		Age:  15,
		Sex:  "男",
	})
	fmt.Printf("%v\n", lruCache.Get("A"))
}
