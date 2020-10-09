/*
 *@Description
 *@author          lirui
 *@create          2020-09-28 19:41
 */
package main

import (
	"awesomego/onemod/k8s/broadcastertest/kubeconfig"
	"bytes"
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func main() {
	log := make(chan string)

	podLogOpts := corev1.PodLogOptions{
		Container: "nginx",
	}

	// creates the clientset
	clientset := kubeconfig.GetKube()

	req := clientset.CoreV1().Pods("default").GetLogs("nginx-deployment-7cf7cdffcd-d94xq", &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		fmt.Println(err)
	}
	str := buf.String()

	log <- str

	for {
		select {
		case ll := <-log:
			fmt.Println(ll)
		}
	}

}
