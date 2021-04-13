package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var up = websocket.Upgrader{}

func main() {
	http.HandleFunc("/ws", func(responseWriter http.ResponseWriter, request *http.Request) {

		up.CheckOrigin = func(r *http.Request) bool { return true }

		wcconn, err := up.Upgrade(responseWriter, request, nil)

		if err != nil {
			fmt.Printf("up falid %v", err)
		}

		defer wcconn.Close()

		for {

			// 读消息
			mt, message, err := wcconn.ReadMessage()

			if err != nil {
				fmt.Printf("read message faile %v", err)
				break
			}

			fmt.Printf("消息类型为：%v,消息内容为：%v \n", mt, string(message))

			// 写消息
			err = wcconn.WriteMessage(mt, []byte("hhhhh"))
			if err != nil {
				fmt.Printf("写消息错误", err)
				break
			}
		}

	})

	http.ListenAndServe(":1989", nil)
}
