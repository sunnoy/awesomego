/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 19:46
 */
package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map

	sm.Store("aa", "kkkk")
	sm.Store("sss", "iiii")

	fmt.Println(sm.Load("aas"))

	sm.Range(func(key, value interface{}) bool {
		fmt.Println(key, value)
		return true
	})
}
