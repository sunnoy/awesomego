/*
 *@Description
 *@author          lirui
 *@create          2021-04-13 19:25
 */
package main

import (
	"fmt"
	"net"
)

func main() {
	var d net.Dialer
	con, err := d.Dial("tcp", "127.0.0.1:1989")
	if err != nil {
		fmt.Printf("dial fiald %v", err)
	}

	_, err = con.Write([]byte("i am sss"))
	if err != nil {
		return
	}

	buf := make([]byte, 4096)
	cnt, err := con.Read(buf)
	if err != nil || cnt == 0 {
		con.Close()
	}
	fmt.Printf("收到的内容是 %v", string(buf[0:cnt]))

}
