package main

// 定义了一个接口
type say interface {
	// 包含了一个方法
	spack()
}

// 创建了一个结构体
type dong struct {
	name string
}

// 接头体的方法
// 给用户定义的类型添加行为
func (d dong) spack() {
	print(d.name)
}

func main() {
	dd := dong{
		name: "ssss",
	}

	notsay(dd)
}

func notsay(s say) {
	s.spack()
}
