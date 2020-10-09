package main

import (
	"fmt"
	"reflect"
)

type User struct {
	Name string
	Age  int
}

func main() {
	u := User{
		Name: "sss",
		Age:  88,
	}

	// %T 和 reflect.TypeOf() 一样
	fmt.Println(reflect.TypeOf(u))
	fmt.Println()
	fmt.Printf("%T", u)
	fmt.Println()
	fmt.Println()
	// %v 和 reflect.ValueOf() 一样
	fmt.Printf("%v", u)
	fmt.Println()
	fmt.Println(reflect.ValueOf(u))
}
