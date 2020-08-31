package controller

import (
	"awesomego/ginscket/server"
	"github.com/gin-gonic/gin"
)

type NameController struct {
}

func NewNameController() *NameController {
	return &NameController{}
}

func (namecontroller *NameController) Router(server *server.Sever) {
	server.Handle("GET", "/name", namecontroller.GetName())
}

func (namecontroller *NameController) GetName() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(200, gin.H{
			"data": "dongdong",
		})
	}
}
