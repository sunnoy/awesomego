/*
 *@Description
 *@author          lirui
 *@create          2020-09-25 16:16
 */
package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", SerHome)
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
	}

	// 完成协议升级
	// ws 是 The Conn type represents a WebSocket connection.
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go WriteWs(ws)

}

func WriteWs(ws *websocket.Conn) {
	defer ws.Close()

	// 启动两个技计时器，用于计时心跳和发送消息
	pingTicker := time.NewTicker(2 * time.Second)
	fileTicker := time.NewTicker(1 * time.Second)

	for {
		select {

		case <-pingTicker.C:

		case <-fileTicker.C:

			if err := ws.WriteMessage(websocket.TextMessage, []byte("aaa")); err != nil {
				return
			}

		}

	}

}

var homeTempl = template.Must(template.New("name").Parse(homeHTML))

// 构造响应来发网客户端
// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.

// 代表了从客户端过来的请求实体
// A Request represents an HTTP request received by a server
// or to be sent by a client.
func SerHome(w http.ResponseWriter, r *http.Request) {
	var v = struct {
		Host    string
		Data    string
		LastMod string
	}{
		r.Host,
		"aaaa",
		"aaaa",
	}

	// Execute applies a parsed template to the specified data object,
	// writing the output to wr.
	homeTempl.Execute(w, v)
}

const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log(evt.data);
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>
`
