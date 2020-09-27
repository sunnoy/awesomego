/*
 *@Description
 *@author          lirui
 *@create          2020-09-26 15:39
 */
package main

func main() {

	ni := struct { //声明结构体
		name string
		age  int
		int  // 匿名字段 同一种只能使用一个
	}{ // 实例化结构体
		name: "ss",
		age:  8,
		int:  8, // 用类型名称来当做字段值
	}

	println(ni)

}
