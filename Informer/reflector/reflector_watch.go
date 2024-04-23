package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

// create pods list & watch
func main() {
	// helper 只是一个类似上文演示的 config, 只要用于初始化各种客户端
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	cliset, err := kubernetes.NewForConfig(config)
	lwc := cache.NewListWatchFromClient(cliset.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	watcher, err := lwc.Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case v, ok := <-watcher.ResultChan():
			if ok {
				fmt.Println(v.Type, ":", v.Object.(*v1.Pod).Name, "-", v.Object.(*v1.Pod).Status.Phase)
			}

		}
	}
}
