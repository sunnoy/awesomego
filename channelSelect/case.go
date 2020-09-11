package main

import (
	"fmt"
	"time"
)

func main() {
	c1 := make(chan string, 1) //定义两个有缓冲通道，容量分别为1
	c2 := make(chan string, 1)

	go func() { //定义一个协程
		time.Sleep(time.Second * 4) //隔1秒发送数据
		c1 <- "name: xuchao"        //向c1通道发送数据
	}()

	go func() {
		time.Sleep(time.Second * 6) //隔6秒发送数据
		c2 <- "age: 25"             //向c2通道发送数据
	}()

	//for i := 0; i < 2; i++ {                  //使用select来获取这两个通道的值，然后输出
	for { //使用select来获取这两个通道的值，然后输出
		tm := time.NewTimer(time.Second * 1) //给通道创建容忍时间，如果5s内无法读写，就即刻返回

		select {
		case msg1 := <-c1: //接收c1通道数据（消费数据）
			fmt.Println(msg1)
		case msg2 := <-c2: //接收c2通道数据（消费数据）
			fmt.Println(msg2)
		case <-tm.C:
			fmt.Println("send data timeout!")
		}
	}
}
