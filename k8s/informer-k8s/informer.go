package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	// 创建一个通道
	stopper := make(chan struct{})

	// 最后关闭这个通道
	defer close(stopper)

	// factory 是一个 SharedInformerFactory 接口，
	// SharedInformerFactory provides shared informers for resources in all known
	// API group versions
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// 指定要监听的资源对象
	podInFormer := factory.Core().V1().Pods()

	// 每种资源对象都实现了一个接口
	// 	Informer() cache.SharedIndexInformer
	//	Lister() v1.podLister
	informer := podInFormer.Informer()

	defer runtime.HandleCrash()

	// 实例化后需要启动 informer
	// list 并且 watch
	go factory.Start(stopper)

	// 检测是否同步完成
	// WaitForCacheSync waits for caches to populate.  It returns true if it was successful

	// HasSynced returns true if the shared informer's store has been
	// informed by at least one full LIST of the authoritative state
	// of the informer's object collection.  This is unrelated to "resync".
	isSynced := cache.WaitForCacheSync(stopper, informer.HasSynced)

	if isSynced == false {
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
		return
	}

	// 创建自定的资源事件处理函数
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addfun,
		UpdateFunc: func(interface{}, interface{}) { fmt.Println("update not implemented") }, // 此处省略 workqueue 的使用
		DeleteFunc: func(interface{}) { fmt.Println("delete not implemented") },
	})

	// 获取list之后的对象
	podLister := podInFormer.Lister()

	podList, err := podLister.List(labels.Everything())

	fmt.Println("podlist", podList[0])

	<-stopper

}

func addfun(obj interface{}) {

	// 类型断言，如果obj类型是括号后的类型则断言成功
	// 第一个参数就是括号内的类型
	// 第二个参数布尔值 表示是否断言成功
	//pod,ok := obj.(*v1.Pod)
	pod := obj.(*v1.Pod)

	println(pod.Name)
}
