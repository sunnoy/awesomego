package controller

import (
	"awesomego/gin/ginscket/server"
	"github.com/gin-gonic/gin"
)

type UserController struct {
}

// 返回结构体 UserController
// 这个结构体的需要实现接口 IController
// 也就是说结构体 UserController 是 IController 类型的
func NewUserController() *UserController {
	return &UserController{}
}

func (UserController *UserController) Router(server *server.Sever) {
	server.Handle("GET", "/", UserController.GetUser())
}

func (UserController *UserController) GetUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{
			"data": "hello world",
		})
	}
}
