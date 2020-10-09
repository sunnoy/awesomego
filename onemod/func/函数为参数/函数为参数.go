package main

func main() {
	a, b := 3, 5

	println(add(a, b, func(i int, i2 int) int {
		return i + i2
	}))
}

func add(a, b int, fun func(int, int) int) int {
	return fun(a, b)

}
