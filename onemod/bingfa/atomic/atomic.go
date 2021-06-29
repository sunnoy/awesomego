/*
 *@Description
 *@author          lirui
 *@create          2021-06-28 19:20
 */
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	NodeName string
	Addr     string
	Count    int32
}

func loadNewConfig() Config {
	return Config{
		NodeName: "北京",
		Addr:     "10.77.95.27",
		Count:    rand.Int31(),
	}
}

func main() {
	// 示例化value类型
	var config atomic.Value

	// 存入新的配置
	config.Store(loadNewConfig())

	var cond = sync.NewCond(&sync.Mutex{})

	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Second)
		}

		config.Store(loadNewConfig())

		cond.Broadcast() // 通知等待着配置已变更

	}()

	go func() {
		for {
			cond.L.Lock()
			cond.Wait() // 等待Broadcast()方法/*

			c := config.Load().(Config) // 读取新的配置
			fmt.Printf("new config: %+v\n", c)

			cond.L.Unlock()
		}
	}()

	select {}
}
