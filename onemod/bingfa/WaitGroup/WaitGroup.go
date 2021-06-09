/*
 *@Description
 *@author          lirui
 *@create          2021-06-09 19:41
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	count uint64
}

func (c *Counter) Incr() {
	c.mu.Lock()
	c.count++
	c.mu.Unlock()
}

func (c *Counter) Count() uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func worker(c *Counter, wg *sync.WaitGroup) {
	// 每走完了一个worker线程就通过调用done方法将wg计数器减1
	defer wg.Done()
	time.Sleep(time.Second)
	c.Incr()
}

func main() {
	var counter Counter
	var wg sync.WaitGroup

	// 初始化wg的值为10
	wg.Add(10)

	// 创建10个goroutine
	for i := 0; i < 10; i++ {
		go worker(&counter, &wg)

	}

	// 阻塞当前，主goroutine等待wg变为0
	wg.Wait()

	fmt.Println(counter.Count())
}
