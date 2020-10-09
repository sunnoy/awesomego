/*
 *@Description
 *@author          lirui
 *@create          2020-10-05 21:08
 */
package main

import (
	"github.com/gin-gonic/gin"
)

//var clinetset = kubeconfig.GetKube()

func main() {
	r := gin.Default()
	r.GET("/api/v1/namespaces/:namespace/pods/:name/log", func(context *gin.Context) {
		url := context.Request.URL.String()
		context.JSON(200, gin.H{
			"url":       url,
			"namespace": context.Param("namespace"),
			"pod-name":  context.Param("name"),
		})
	})
	r.Run(":1989")
}
