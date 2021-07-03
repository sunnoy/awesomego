/*
 *@Description
 *@author          lirui
 *@create          2020-09-25 17:50
 */
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	http.HandleFunc("/ws", echo)
	//http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":1989", nil))

}

func echo(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer ws.Close()

}
