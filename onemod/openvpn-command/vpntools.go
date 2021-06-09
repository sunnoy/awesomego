/*
 *@Description
 *@author          lirui
 *@create          2021-05-17 18:52
 */
package main

import (
	expect "github.com/Netflix/go-expect"
	"log"
	"os/exec"
	"time"
)

func main() {

	c, err := expect.NewConsole()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command("sh", "-c", "read name1;echo name1is$name1;read name2;echo name2is$name2")

	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()
	cmd.Stdin = c.Tty()

	go func() {
		c.ExpectEOF()
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal("err is", err)
	}

	time.Sleep(time.Second)

	c.ExpectString("name1")

	c.Send("name1111")

}
