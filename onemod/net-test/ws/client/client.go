package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	url2 "net/url"
	"time"
)

func main() {
	url := url2.URL{
		Scheme: "ws",
		Host:   "127.0.0.1:1989",
		Path:   "/ws",
	}

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		fmt.Printf("dial falie %v", err)
	}

	defer c.Close()

	// 读消息
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read err", err)
				break
			} else {
				log.Printf("recv is %s", message)
			}
		}
	}()

	// 写消息
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write err", err)
			} else {
				log.Println(t.String())
			}

		}

	}

}
