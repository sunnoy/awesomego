/*
 *@Description
 *@author          lirui
 *@create          2020-10-14 14:38
 */
package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:1990")
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()

	status := make(chan string)
	go Write(conn, status)
	go Read(conn, status)

	fmt.Println(<-status)
	fmt.Println(<-status)

}

func Write(conn net.Conn, status chan string) {
	for i := 0; i < 10; i++ {
		wn, err := conn.Write([]byte("hello " + strconv.Itoa(i) + "\r\n"))
		if err != nil {
			log.Println(err)
		}
		log.Println(wn)
	}

	status <- "write done"
}

func Read(conn net.Conn, status chan string) {

	buf := make([]byte, 1024)

	rr, err := conn.Read(buf)

	if err != nil {
		log.Println(err)
	}
	log.Println(rr)

	status <- "read done"
}
