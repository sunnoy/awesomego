/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 15:09
 */
package main

import (
	"awesomego/k8s/broadcastertest/broadcaster"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"log"
	"net/http"
	"time"
)

// 初始化 SocketData 通道
var ResourcesChannel = make(chan broadcaster.SocketData)

var ResourcesBroadcaster = broadcaster.NewBroadcaster(ResourcesChannel)

func main() {

	r := gin.Default()

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
	go readControl(connection, b, subscriber)
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
	if err := connection.WriteMessage(websocket.TextMessage, payload); err != nil {
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

func NewPodController(kind string, messagetype broadcaster.MessageType) {

	//NewController()
}

func NewController(kind string, informer cache.SharedIndexInformer, onCreated, onUpdated, onDeleted broadcaster.MessageType, filter func(interface{}, bool) interface{}) {
	println("In NewController")

	if filter == nil {
		filter = func(obj interface{}, skipDeletedCheck bool) interface{} {
			return obj
		}
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			println("Controller detected %s '%s' created", kind, obj.(metav1.Object).GetName())
			data := broadcaster.SocketData{
				MessageType: onCreated,
				Payload:     filter(obj, true),
			}
			ResourcesChannel <- data
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldSecret, newSecret := oldObj.(metav1.Object), newObj.(metav1.Object)
			// If resourceVersion differs between old and new, an actual update event was observed
			if oldSecret.GetResourceVersion() != newSecret.GetResourceVersion() {
				println("Controller detected %s '%s' updated", kind, oldSecret.GetName())
				data := broadcaster.SocketData{
					MessageType: onUpdated,
					Payload:     filter(newObj, true),
				}
				ResourcesChannel <- data
			}
		},
		DeleteFunc: func(obj interface{}) {
			println("Controller detected %s '%s' deleted", kind, GetDeletedObjectMeta(obj).GetName())
			data := broadcaster.SocketData{
				MessageType: onDeleted,
				Payload:     filter(obj, false),
			}
			ResourcesChannel <- data
		},
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
