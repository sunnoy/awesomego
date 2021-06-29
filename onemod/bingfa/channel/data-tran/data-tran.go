/*
 *@Description
 *@author          lirui
 *@create          2021-06-29 19:49
 */
package main

import (
	"fmt"
	"time"
)

// 定义一个令牌类型
type Token struct {
}

// 读取令牌
func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {

		// 取得了令牌 立即打印出自己的编号
		token := <-ch
		fmt.Println((id + 1))

		time.Sleep(time.Second)

		// 将令牌交给下一家
		nextCh <- token
	}
}

func main() {
	// 创建一个通道切片，每个通道里面放的是一个类型为 Token 的结构体
	chs := []chan Token{
		make(chan Token),
		make(chan Token),
		make(chan Token),
		make(chan Token),
	}

	for i := 0; i < 4; i++ {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	chs[0] <- struct{}{}

	select {}

}
