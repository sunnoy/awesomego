/*
 *@Description
 *@author          lirui
 *@create          2020-09-25 16:16
 */
package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func main() {

	http.HandleFunc("/ws", SerWS)

	if err := http.ListenAndServe(":1989", nil); err != nil {
		log.Fatal(err)
	}
}

func SerWS(w http.ResponseWriter, r *http.Request) {
	// 指定协议升级时候的参数
	// Upgrader specifies parameters for upgrading an HTTP connection to a
	// WebSocket connection.
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 不校验host头部
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 完成协议升级
	// wspkg 是 The Conn type represents a WebSocket connection.
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	ver := r.URL.Query().Get("ver")

	go WriteWs(ws, ver)

}

func WriteWs(ws *websocket.Conn, ver string) {
	defer ws.Close()

	// 启动两个技计时器，用于计时心跳和发送消息
	pingTicker := time.NewTicker(2 * time.Second)
	fileTicker := time.NewTicker(1 * time.Second)

	for {
		select {

		case <-pingTicker.C:

		case <-fileTicker.C:

			if err := ws.WriteMessage(websocket.TextMessage, []byte(ver)); err != nil {
				return
			} else {
				fmt.Println("print ", ver)
			}

		}

	}

}

/**
	curl --include \
     --no-buffer \
     --header "Connection: Upgrade" \
     --header "Upgrade: websocket" \
     --header "Sec-WebSocket-Key: SGVsbG8sIHdvcmxkIQ==" \
     --header "Sec-WebSocket-Version: 13" \
    http://127.0.0.1:1989/ws?ver=v2
**/
