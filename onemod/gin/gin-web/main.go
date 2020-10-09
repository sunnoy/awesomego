package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/a", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run("0.0.0.0:9090")
}
