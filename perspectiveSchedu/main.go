package main

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	clientSet := GetConfig()

	nodes, err := clientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		println(err)
	}

	for _, node := range nodes.Items {

		pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + node.Name,
		})

		if err != nil {
			println(err)
		}

		println(len(pods.Items))

		//
		//
		//
		//fmt.Printf("Node[ %v ]（nodeport）上面的业务pod有：\n", node.Name)
		//for _, pod := range pods.Items {
		//	if _, ok := pod.Labels["env"]; ok {
		//		fmt.Println("pod: ", pod.Name)
		//	}
		//}

	}

}
