/*
 *@Description
 *@author          lirui
 *@create          2020-10-09 11:46
 */
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	server := &APIServer{
		engine: gin.Default(),
	}
	server.registryApi()
	server.engine.Run(":38080")
}

type APIServer struct {
	engine *gin.Engine
}

func (s *APIServer) registryApi() {
	registryStream(s.engine)
}

func registryStream(engine *gin.Engine) {
	engine.GET("/stream", func(ctx *gin.Context) {
		// http响应对象
		w := ctx.Writer
		header := w.Header()
		//在响应头添加分块传输的头字段Transfer-Encoding: chunked
		header.Set("Transfer-Encoding", "chunked")
		//header.Set("Content-Type", "text/html")
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		for i := 0; i < 10; i++ {
			w.Write([]byte(fmt.Sprintf("%d", i)))
			w.(http.Flusher).Flush()
			time.Sleep(time.Duration(1) * time.Second)
		}

	})
}

// curl 请求 -N 不需要buffer
// curl -v 127.0.0.1:38080/stream -N
