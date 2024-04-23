package main

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	// get specifc pod
	pod, err := clientset.CoreV1().Pods("kube-system").Get(context.TODO(), "etcd-minikube", v1.GetOptions{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("pod NAME %v", pod.Name)
	}

	// get pod list

}
