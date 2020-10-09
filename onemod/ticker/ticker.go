package main

import (
	"fmt"
	"time"
)

func main() {

	//myTicker := time.NewTicker(time.Second) //设置时间周期
	//for {
	//	nowTime := <-myTicker.C //当前时间
	//	if nowTime.Hour() == 11 && nowTime.Minute() == 35 {
	//		fmt.Println("Golang")
	//		break
	//	}
	//}

	// 定义一个定时器 周期性的向通道发送消息
	ticker := time.NewTicker(time.Second * 3) // 每隔1s进行一次打印
	for {
		// 会将当前的时间通过通道发送出来
		// 2020-09-09 20:34:08.359085 +0800 CST m=+4.003485721
		//<-ticker.C

		// 第一次接收
		fmt.Println(<-ticker.C)
		// 下一次接收时间就会延长 time.Second*3
		fmt.Println(<-ticker.C)
		fmt.Println("这是ticker的打印")
		// 好的
	}
}
