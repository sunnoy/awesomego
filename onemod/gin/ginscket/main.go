package main

import (
	"awesomego/onemod/gin/ginscket/controller"
	"awesomego/onemod/gin/ginscket/server"
)

func main() {
	server.Init().
		// route 方法参数是个 IController 类型
		Route(
			controller.NewUserController(),
			controller.NewNameController(),
		).
		//GroupRouter(
		//	"v1",
		//	controller.NewUserController()).
		Listen()
}
