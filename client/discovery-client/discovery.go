package main

import (
	"fmt"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	discover, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	resources, ApiResourceLists, err := discover.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	//fmt.Println("resources: %v", resources)
	for i, list := range resources {
		fmt.Printf("index %v \t NAME: %v \t Version %v \t APIVersion %v Kind %+v\n", i, list.Name, list.Versions, list.APIVersion, list.Kind)
	}

	for i, list := range ApiResourceLists {
		for j, resource := range list.APIResources {
			fmt.Printf("list %v\t index %v\t GroupVersion %v\t NAME: %v\t  Verb %+v\n", i, j, list.GroupVersion, resource.Name, resource.Verbs)
		}
	}
}
