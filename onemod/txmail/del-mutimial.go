/*
 *@Description
 *@author          lirui
 *@create          2021-06-24 15:19
 */
package main

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "admin@test.cn")
	m.SetHeader("To", "wode@qq.com")
	m.SetHeader("Subject", "go语言发送的邮件")
	m.SetBody("text/html", "<span style='color:red;'>test</span>")
	//m.Attach("/Users/zc/Pictures/skip.gif")
	d := gomail.NewDialer("smtp.exmail.qq.com", 987, "lirui@xylink.com", "tBF5BuSk6tudDiqX")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
