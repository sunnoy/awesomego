/*
 *@Description
 *@author          lirui
 *@create          2021-06-08 19:40
 */
package main

import (
	"sync"
	"time"
)

type Counter struct {
	mu    sync.RWMutex
	count uint64
}

// 加1封装
func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 返回结构封装，在这里不加锁也行
func (c *Counter) Count() uint64 {
	c.mu.RLocker()

	// 等待追后的结构返回后进行解锁
	defer c.mu.RUnlock()

	return c.count
}

func main() {

	var counter Counter

	for i := 0; i < 10; i++ {
		go func() {

			for {
				counter.Count()
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for {
		counter.Incr()
		time.Sleep(time.Second)
	}

	println(counter.Count())
}
