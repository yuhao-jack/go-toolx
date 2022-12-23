package ring_buf

import (
	"fmt"
	"github.com/yuhao-jack/go-toolx/algorithm/cache/ring_buf"
	"sync"
	"testing"
	"time"
)

func TestRingBuf(t *testing.T) {
	ringBufCache := ring_buf.NewRingBufCache[int](3)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for i := 0; i < 10; i++ {
			if ringBufCache.Put(i + 5) {
				fmt.Println("put ", i+5, " success")
			} else {
				fmt.Println("put ", i+5, " failed")
			}
			time.Sleep(time.Second)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10; i++ {
			get, b := ringBufCache.Get(-1)
			fmt.Println("get:", get, b)
			time.Sleep(2 * time.Second)
		}
		wg.Done()
	}()

	wg.Wait()

	fmt.Println(ringBufCache.GetAll())

}
