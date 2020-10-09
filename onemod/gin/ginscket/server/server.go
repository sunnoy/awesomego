package server

import (
	"github.com/gin-gonic/gin"
)

// 定义一个接口
type IController interface {
	// 包含一个 router 方法
	Router(server *Sever)
}

type Sever struct {
	*gin.Engine
	g *gin.RouterGroup
}

// 初始化一个 gin Engine 实例
// 这个实例包含在 server 结构体中
func Init() *Sever {
	s := &Sever{
		Engine: gin.New(),
	}

	return s
}

// 指定监听
func (server *Sever) Listen() {
	server.Run("0.0.0.0:1989")
}

// route 方法用来注册控制器
// 接收一个可变参数 参数类型是接口类型 IController
// 任何实现了这个接口的结构体都可以作为 Route 的参数
func (server *Sever) Route(controllers ...IController) *Sever {
	for _, c := range controllers {
		c.Router(server)
	}

	return server
}

// 注册路由组以及相关的控制器
// 控制器为可变参数
func (server *Sever) GroupRouter(group string, controllers ...IController) *Sever {
	// 将 route group 初始化
	server.g = server.Group(group)

	// 迭代控制器
	for _, c := range controllers {
		c.Router(server)
	}
	return server
}
