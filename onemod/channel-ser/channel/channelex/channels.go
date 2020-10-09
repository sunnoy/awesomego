package main

func main() {

	//runtime.GOMAXPROCS(1)
	ch := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i

		}
		close(ch)
	}()

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
