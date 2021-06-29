/*
 *@Description
 *@author          lirui
 *@create          2021-06-24 19:28
 */
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 初始化根 context
	ctx := context.Background()

	// cancelCtx ---- Copy of the parentContext with the new done channel.
	// cancelFunc ---- A cancel function which when called closes this done channel

	// WithCancel returns a copy of parent with a new Done channel. The returned
	// context's Done channel is closed when the returned cancel function is called
	// or when the parent context's Done channel is closed, whichever happens first.
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	go task(cancelCtx)

	// 假装进行业务处理
	time.Sleep(time.Second * 3)

	// 处理完了调用 取消函数，取消这个done通道
	cancelFunc()
	time.Sleep(time.Second * 1)
}

func task(ctx context.Context) {
	i := 1
	for {
		select {

		// 通道还在，进行阻塞
		// 调用 cancelFunc() 之后，通道
		case <-ctx.Done():
			fmt.Println("Gracefully exit")

			// If Done is not yet closed, Err returns nil.
			// If Done is closed, Err returns a non-nil error explaining why:
			// Canceled if the context was canceled
			// or DeadlineExceeded if the context's deadline passed.
			fmt.Println(ctx.Err(), "取消了") // 返回 context canceled
			return
		default:
			fmt.Println(i)
			time.Sleep(time.Second * 1)
			i++
		}
	}
}
