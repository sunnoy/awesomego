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
	// 写的时候加入写锁
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

// 返回结构封装，在这里不加锁也行
func (c *Counter) Count() uint64 {
	//读的时候加入读锁
	c.mu.RLock()

	// 等待追后的结构返回后进行解锁
	defer c.mu.RUnlock()

	return c.count
}

func main() {

	var counter Counter

	//可以明确区分 reader 和 writer goroutine 的场景，且有⼤量的并发读、少量的并 发写，并且有强烈的性能需求
	for i := 0; i < 10; i++ {
		go func() {

			for {
				println(counter.Count())
				time.Sleep(time.Millisecond)
			}
		}()
	}

	for {
		counter.Incr()
		time.Sleep(time.Second)
	}

}
