package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		var sum int = 0
		for i := 0; i < 10; i++ {
			sum += i
		}
		ch <- sum
		ch <- sum
	}()

	fmt.Println(<-ch)

	time.Sleep(time.Second * 2)

	fmt.Println(<-ch)

}
