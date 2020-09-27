/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 15:09
 */
package main

import (
	"awesomego/k8s/broadcastertest/broadcaster"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"time"
)

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
                var conn = new WebSocket("ws://{{.Host}}/ws");
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

var homeTempl = template.Must(template.New("name").Parse(homeHTML))

// 初始化 SocketData 通道
var ResourcesChannel = make(chan broadcaster.SocketData)

var ResourcesBroadcaster = broadcaster.NewBroadcaster(ResourcesChannel)

func main() {

	SockDate()

	r := gin.Default()

	r.GET("/", func(context *gin.Context) {
		var v = struct {
			Host    string
			Data    string
			LastMod string
		}{
			context.Request.Host,
			"aaaa",
			"aaaa",
		}

		// Execute applies a parsed template to the specified data object,
		// writing the output to wr.
		homeTempl.Execute(context.Writer, v)
	})

	r.GET("/ws", func(c *gin.Context) {
		wscon := SerWS(c.Writer, c.Request)
		WriteOnlyWebsocket(wscon, ResourcesBroadcaster)

	})

	r.Run(":1989")
}

// WriteOnlyWebsocket discards text messages from the peer connection
func WriteOnlyWebsocket(connection *websocket.Conn, b *broadcaster.Broadcaster) {
	// The underlying connection is never closed so this cannot error
	subscriber, _ := b.Subscribe()
	//readControl(connection, b, subscriber)
	write(connection, subscriber)

}

// ping over the socket with a given deadline; if there's an error, close
func writePing(connection *websocket.Conn, deadline time.Time) {
	if err := connection.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
		ReportClosing(connection)
	}
}

// readControl will unsubscribe on connection failures
func readControl(connection *websocket.Conn, b *broadcaster.Broadcaster, s *broadcaster.Subscriber) {
	// Connection lifecycle handler
	connection.SetPongHandler(func(string) error {
		// Extend deadline to prevent expiration
		deadline := time.Now().Add(time.Second * 2)
		connection.SetReadDeadline(deadline)
		// Cut down on ping/pong traffic
		time.Sleep(time.Second)
		// Ellicit another ping
		writePing(connection, deadline)
		return nil
	})
	initialDeadline := time.Now().Add(time.Second)
	connection.SetReadDeadline(initialDeadline)
	// Kick off cycle
	writePing(connection, initialDeadline)
	for {
		// Connection has either decayed or close has been requested from server side
		if _, _, err := connection.ReadMessage(); err != nil {
			//logging.Log.Error("websocket connection to client lost: ", err)
			b.Unsubscribe(s)
			return
		}
	}
}

// ReportClosing sends close to client then closes connection
func ReportClosing(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	connection.Close()
}

// Send data over the connection using the subscriber channel, if there's a failure we return
func write(connection *websocket.Conn, subscriber *broadcaster.Subscriber) {
	subChan := subscriber.SubChan()
	unsubChan := subscriber.UnsubChan()
	for {
		select {
		case socketData := <-subChan:
			if !websocketSend(connection, socketData) {
				return
			}
		case <-unsubChan:
			return
		}
	}
}

// Returns whether successful or not, closes connection on failures
func websocketSend(connection *websocket.Conn, data broadcaster.SocketData) bool {
	payload, err := json.Marshal(data)
	if err != nil {
		//logging.Log.Errorf("failed to marshal status: %s", err)
		ReportClosing(connection)
		return false
	}

	println(payload)

	payloads := []byte("sssss")
	if err := connection.WriteMessage(websocket.TextMessage, payloads); err != nil {
		//logging.Log.Errorf("could not write the message to the websocket client connection, error: %s", err)
		ReportClosing(connection)
		return false
	}
	return true
}

func SerWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
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
		return nil
	}

	return ws

}

func SockDate() {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "/Users/lirui/.kube/config6")
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)

	if err != nil {
		fmt.Println(err)
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, 0)

	PodInformer := informerFactory.Core().V1().Pods()

	PodInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("没有定义更新")
			data := broadcaster.SocketData{
				MessageType: broadcaster.NamespaceCreated,
				Payload:     oldObj,
			}

			ResourcesChannel <- data

		},
		DeleteFunc: nil,
	})

}

func GetDeletedObjectMeta(obj interface{}) metav1.Object {
	// Deal with tombstone events by pulling the object out.  Tombstone events wrap the object in a
	// DeleteFinalStateUnknown struct, so the object needs to be pulled out.
	// Copied from sample-controller
	// This should only happen when we're missing events.
	if _, ok := obj.(metav1.Object); !ok {
		// If the object doesn't have Metadata, assume it is a tombstone object of type DeletedFinalStateUnknown
		if tombstone, ok := obj.(cache.DeletedFinalStateUnknown); !ok {
			println("Error decoding object: Expected cache.DeletedFinalStateUnknown, got %T", obj)
			return &metav1.ObjectMeta{}
		} else {
			// Set obj to the tombstone obj
			obj = tombstone.Obj
		}
	}

	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err != nil {
		println("Missing meta for object %T: %v", obj, err)
		return &metav1.ObjectMeta{}
	} else {
		return o
	}
}
