/*
 *@Description
 *@author          lirui
 *@create          2020-10-14 14:32
 */
package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {

	//listen:= net.ListenConfig{
	//	KeepAlive: 4*time.Second,
	//}
	//
	//listen.Listen(context.TODO(),"","")

	l, err := net.Listen("tcp", "127.0.0.1:1990")

	if err != nil {
		log.Println(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		if err != nil {
			log.Println(err)
		}

		io.Copy(conn, conn)

	}
}
