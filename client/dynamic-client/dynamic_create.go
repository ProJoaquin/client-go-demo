package main

import (
	"context"
	_ "embed"
	"github.com/ProJoaquin/client-go-demo/client/config"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
)

// 这个是新特性使用注释加载配置
//
//go:embed deployment.yaml
var deployTpl string

// dynamic client 创建 Deploy
func main() {
	k8sConfig := config.NewK8sConfig()
	// 动态客户端
	dynamicCli := k8sConfig.InitDynamicClient()

	// 可以随意指定集群拥有的资源, 进行创建
	deployGVR := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	deployObj := &unstructured.Unstructured{}
	if err := yaml.Unmarshal([]byte(deployTpl), deployObj); err != nil {
		log.Fatalln(err)
	}

	if _, err := dynamicCli.
		Resource(deployGVR).
		Namespace("default").
		Create(context.Background(), deployObj, metav1.CreateOptions{}); err != nil {
		log.Fatalln(err)
	}

	log.Println("Create deploy succeed")
}
