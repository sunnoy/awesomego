package main

import "time"

func send(ch chan int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
}
func receive(ch chan int) {
	//for {
	//	fmt.Println(<-ch)
	//}

	for c := range ch {
		println(c)
	}
	close(ch)
}
func main() {
	ch := make(chan int)
	go send(ch)
	go receive(ch)
	time.Sleep(2 * time.Millisecond)
}
