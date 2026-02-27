package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Inventory struct {
	stock   int
	rwMutex sync.RWMutex
}

func (i *Inventory) getStock() int {
	i.rwMutex.RLock()
	defer i.rwMutex.RUnlock()
	return i.stock
}

func (i *Inventory) deductStock(quantity int) bool {
	i.rwMutex.Lock()
	defer i.rwMutex.Unlock()
	if i.stock-quantity < 0 {
		fmt.Println("库存不足")
		return false
	}
	time.Sleep(100 * time.Millisecond)
	i.stock = i.stock - quantity
	return true
}

// Counter 计数器, 使用原子操作
// 原子操作会比加锁快一个或多个数量级
type Counter struct {
	value int64
}

func (c *Counter) Int() {
	time.Sleep(100 * time.Millisecond)
	atomic.AddInt64(&c.value, 1)
}
func (c *Counter) Dec() {
	atomic.AddInt64(&c.value, -1)
}
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

func main() {
	// * 模拟扣减库存
	inventory := &Inventory{
		stock: 100,
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inventory.deductStock(1)
		}()
	}

	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("剩余库存: ", inventory.getStock())
	}

	wg.Wait()
	fmt.Println("最终剩余库存: ", inventory.getStock())

	// * 模拟点赞
	var c Counter
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Int()
		}()
	}

	// 取消点赞
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Dec()
		}()
	}

	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("点赞数: ", c.Get())
	}

	wg.Wait()
	fmt.Println("最终点赞数:", c.Get()) // 90
}
