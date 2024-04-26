package main

import (
	"github.com/ProJoaquin/client-go-demo/pkg"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func main() {
	// 1. config
	// 2. clientset
	// 3. informer
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalf("pull config failed")
		}
		config = inClusterConfig
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("create clientset failed")
	}

	//
	factory := informers.NewSharedInformerFactory(clientset, 0)
	servicesInformer := factory.Core().V1().Services()
	ingressInformer := factory.Networking().V1().Ingresses()
	controller := pkg.Newcontroller(clientset, servicesInformer, ingressInformer)
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	controller.Run(stopCh)
}
