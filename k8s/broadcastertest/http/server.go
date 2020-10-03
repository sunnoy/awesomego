/*
 *@Description
 *@author          lirui
 *@create          2020-09-27 15:09
 */
package main

import (
	"awesomego/k8s/broadcastertest/broadcaster"
	"awesomego/k8s/broadcastertest/kubeconfig"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"html/template"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"log"
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

func main() {

	stopper := make(chan struct{})

	defer close(stopper)

	factory := informers.NewSharedInformerFactory(kubeconfig.GetKube(), 0)

	podInFormer := factory.Core().V1().Pods()

	informer := podInFormer.Informer()

	go factory.Start(stopper)

	isSynced := cache.WaitForCacheSync(stopper, informer.HasSynced)

	if isSynced == false {
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
		return
	} else {
		fmt.Println("chengogn")
	}

	ResourcesChannel := make(chan broadcaster.SocketData)
	ResourcesBroadcaster := broadcaster.NewBroadcaster(ResourcesChannel)

	podSubscribe, err := ResourcesBroadcaster.Subscribe()
	if err != nil {
		fmt.Println(err)
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("从通道里面塞数据", pod.Name)
			data := broadcaster.SocketData{
				MessageType: "sss",
				Payload:     pod.Name,
			}
			ResourcesChannel <- data

		},
		UpdateFunc: nil,
		DeleteFunc: nil,
	})

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

	r.GET("/ws", func(context *gin.Context) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		wscon, err := upgrader.Upgrade(context.Writer, context.Request, nil)

		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				log.Println(err)
			}
		}

		subChan := podSubscribe.SubChan()

		unsubChan := podSubscribe.UnsubChan()
		for {
			select {
			case socketData := <-subChan:
				payload, err := json.Marshal(socketData)
				if err != nil {
					fmt.Printf("failed to marshal status: %s", err)

				}

				if err := wscon.WriteMessage(websocket.TextMessage, payload); err != nil {
					fmt.Printf("could not write the message to the wspkg client connection, error: %s", err)

				}
			case <-unsubChan:
				return
			}
		}

		//date, _ := json.Marshal(<-podSubscribe.SubChan())
		//
		//if err  = wscon.WriteMessage(websocket.TextMessage, date); err != nil {
		//	fmt.Println(err)
		//}

	})

	r.Run(":1989")

	//
	//for v := range podSubscribe.SubChan() {
	//	fmt.Println("从通道里面取数据")
	//	fmt.Println(v.MessageType)
	//}

	<-stopper

}
