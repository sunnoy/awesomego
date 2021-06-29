/*
 *@Description
 *@author          lirui
 *@create          2021-06-29 19:20
 */
package main

import (
	"fmt"
	"reflect"
)

// 通过反射使用select处理不定数量的chan

func main() {
	var ch1 = make(chan int, 10)
	var ch2 = make(chan int, 10)

	var cases = createCases(ch2, ch1)

	for i := 0; i < 10; i++ {
		// select 语句
		chosen, recy, ok := reflect.Select(cases)
		if recy.IsValid() {
			fmt.Println("recv", cases[chosen].Dir, recy, ok)
		} else {
			fmt.Println("send", cases[chosen].Dir, ok)
		}
	}

}

// 将不定数量的chan处理让后返回一个 reflect.SelectCase 切片
func createCases(chs ...chan int) []reflect.SelectCase {
	var cases []reflect.SelectCase

	// 为每个ch 生成接收的 case
	for _, ch := range chs {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		})
	}

	// 为每个ch生成发送的 case
	for i, ch := range chs {
		v := reflect.ValueOf(i)
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectSend,
			Chan: reflect.ValueOf(ch),
			Send: v,
		})
	}

	return cases

}
