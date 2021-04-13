package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	url2 "net/url"
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

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

}
