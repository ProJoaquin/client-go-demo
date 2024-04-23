package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}
	// RestClient 需要手动设置 1）API Path 请求的HTTP路径；2）资源组和版本；3）NegotiatedSerializer数据的编解码器
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs
	config.APIPath = "api"

	//client
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	// get Pods data
	result := &corev1.PodList{}
	err = restClient.Get().Namespace("kube-system").Resource("pods").VersionedParams(&metav1.ListOptions{Limit: 500}, scheme.ParameterCodec).Do(context.TODO()).Into(result)
	if err != nil {
		panic(err)
	}
	for i, d := range result.Items {
		fmt.Printf("index %v\t NAMESPACE: %v \t NAME %v STATUS %+v\n", i, d.Namespace, d.Name, d.Status.Phase)
	}

	// get specific pod data
	pod := &corev1.Pod{}
	err = restClient.Get().Namespace("kube-system").Resource("pods").Name("etcd-minikube").Do(context.TODO()).Into(pod)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("STATU %v\n", pod.Status.Message)
	}
}
