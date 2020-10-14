package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

func main() {

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	endPoint := "https://www.baidu.com"

	req, err := http.NewRequest("get", endPoint, nil)
	if err != nil {
		log.Fatalf("Error Occured. %+v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpClient.Do(req)
	if err != nil && response == nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	// Let's check if the work actually is done
	// We have seen inconsistencies even when we get 200 OK response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	log.Println("Response Body:", string(body))
}
