package main

import (
	"log"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 创建clientset
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 创建 ratelimitworkqueue

	// 创建informer
	sharedInformers := informers.NewSharedInformerFactory(clientset, time.Minute)
	informer := sharedInformers.Core().V1().Pods().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*v1.Pod)
			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(*v1.Pod)
			nObj := newObj.(*v1.Pod)
			log.Printf("Pod Updated from %s to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*v1.Pod)
			log.Printf("Pod Deleted from Store: %s", mObj.GetName())
		},
	})

	//
	stopCh := make(chan struct{})
	defer close(stopCh)

	informer.Run(stopCh)
}
