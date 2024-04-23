package config

import (
	"k8s.io/client-go/discovery"
	dynamic2 "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type K8sConfig struct {
}

func NewK8sConfig() *K8sConfig {
	return &K8sConfig{}
}

func (c K8sConfig) K8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalln("")
	}
	return config
}

func (c K8sConfig) InitClientset() *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(c.K8sRestConfig())
	if err != nil {
		log.Fatalln("")
	}
	return clientset
}

func (c K8sConfig) InitDynamicClient() dynamic2.Interface {
	dynamic, err := dynamic2.NewForConfig(c.K8sRestConfig())
	if err != nil {
		log.Fatalln("create dynamic failed")
	}
	return dynamic
}

func (c K8sConfig) InitDiscoverClient() *discovery.DiscoveryClient {
	return discovery.NewDiscoveryClient(c.InitClientset().RESTClient())
}
