/*
 *@Description
 *@author          lirui
 *@create          2021-04-13 19:25
 */
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":1989")
	if err != nil {
		fmt.Printf("listen err is %v", err)
	}

	defer l.Close()

	for {
		con, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			buf := make([]byte, 4096)
			cnt, err := c.Read(buf)
			if err != nil || cnt == 0 {
				c.Close()
			}
			fmt.Printf("写的内容是 %v", string(buf[0:cnt]))

			c.Write([]byte("hahahhah"))

			c.Close()
		}(con)
	}

}
