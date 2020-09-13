package main

import "runtime"

func main() {
	runtime.GOMAXPROCS(2)
	ch := make(chan int)

	for i := 0; i < 5; i++ {

		i := i
		go func() {
			ch <- i
			close(ch)
		}()

	}

	for c := range ch {
		println(c)
	}

	//select {
	//case msg := <-ch:
	//	println(msg)
	//default:
	//	println("sss")
	//}

}
