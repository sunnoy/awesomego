/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 19:50
 */
package main

import (
	"fmt"
)

type son struct {
	name string
}

func main() {
	var in interface{}
	in = "aaa"
	m := in.(son)

	// 类型断言，如果obj类型是括号后的类型则断言成功
	// 第一个参数就是括号内的类型
	// 第二个参数布尔值 表示是否断言成功
	//pod,ok := obj.(*v1.Pod)
	//pod := obj.(*v1.Pod)

	fmt.Println(m)
}
