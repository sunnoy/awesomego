package main

import "fmt"

func main() {
	ch := make(chan int)

	go func() {
		var sum int = 0
		for i := 0; i < 10; i++ {
			sum += i
		}
		ch <- sum
	}()

	fmt.Println(<-ch)

}
