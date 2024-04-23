package main

import (
	"fmt"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamic2 "k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	dynamic, err := dynamic2.NewForConfig(config)
	groupVersionResource := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	unstructuredList, err := dynamic.Resource(groupVersionResource).Namespace(corev1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{Limit: 100})
	if err != nil {
		panic(err)
	}
	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredList.UnstructuredContent(), podList)
	if err != nil {
		panic(err)
	}
	for i, d := range podList.Items {
		fmt.Printf("index %v\t NAMESPACE: %v \t NAME %v STATUS %+v\n", i, d.Namespace, d.Name, d.Status.Phase)
	}
}
