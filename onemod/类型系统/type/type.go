package main

import "fmt"

type dong struct {
	name string
	age  int
}

func main() {
	d1 := dong{
		name: "dong",
		age:  4,
	}

	p2 := new(dong)
	p2.name = "sss"

	fmt.Println(d1, *p2)
}
