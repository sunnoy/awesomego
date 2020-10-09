package main

func main() {
	println(get()(2))
}

func get() func(int) int {
	return func(i int) int {
		return i + 3
	}
}
