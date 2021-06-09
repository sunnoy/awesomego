/*
 *@Description
 *@author          lirui
 *@create          2021-06-08 19:40
 */
package main

import "sync"

type Counter struct {
	mu    sync.Mutex
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
	c.mu.Lock()

	// 等待追后的结构返回后进行解锁
	defer c.mu.Unlock()

	return c.count
}

func main() {

	var counter Counter

	var wg sync.WaitGroup

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 100000; i++ {

				counter.Incr()
			}
		}()
	}

	wg.Wait()
	println(counter.Count())
}
