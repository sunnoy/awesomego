/*
 *@Description
 *@author          lirui
 *@create          2020-09-28 15:38
 */
package kubeconfig

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKube() *kubernetes.Clientset {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", "/Users/lirui/Documents/git_repo/xylink/运维日常/kubeconfig/yaml/ops-dev.yaml")
	if err != nil {
		fmt.Println(err)
	}
	clientset, err := kubernetes.NewForConfig(kubeconfig)

	if err != nil {
		fmt.Println(err)
	}
	return clientset
}
