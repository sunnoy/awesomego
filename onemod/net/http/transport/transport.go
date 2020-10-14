/*
 *@Description
 *@author          lirui
 *@create          2020-10-14 11:29
 */
package main

import (
	"log"
	"net/http"
)

type myTransport struct {
}

func main() {

	httpclient := &http.Client{
		Transport: &myTransport{},
	}
	req, _ := http.NewRequest("GET", "/foo", nil)
	res, err := httpclient.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println(res)

}

func (t *myTransport) RoundTrip(rep *http.Request) (*http.Response, error) {
	rep.Header.Add("X-Test", "true")

	return http.DefaultTransport.RoundTrip(rep)
}
