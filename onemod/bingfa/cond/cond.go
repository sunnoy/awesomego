/*
 *@Description
 *@author          lirui
 *@create          2021-06-15 18:51
 */
package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(&sync.Mutex{})

	var ready int

	for i := 0; i < 10; i++ {
		go func(i int) {
			// 模拟运动员准备时间
			time.Sleep(time.Duration(rand.Int63n(10)) * time.Second)

			// ready 属于临界区资源呢 本身就要加锁
			c.L.Lock()

			ready++

			c.L.Unlock()

			log.Printf("运动员#%d 已准备就绪\n", i)

			//用于 条件加锁
			c.Broadcast()

		}(i)
	}

	c.L.Lock()

	for ready != 10 {
		// 用于 条件解锁
		c.Wait()
		log.Println("裁判员被唤醒⼀次:", ready)

	}

	c.L.Unlock()
	log.Println("所有运动员都准备就绪。⽐赛开始，3，2，1, ......")

}
