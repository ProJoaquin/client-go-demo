package main

import (
	"fmt"
	"github.com/ProJoaquin/client-go-demo/client/config"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"time"
)

type CmdHandler struct {
}

// 当接收到添加事件便会执行该回调, 后面的方法以此类推
func (this *CmdHandler) OnAdd(obj interface{}) {
	fmt.Println("Add: ", obj.(*v1.ConfigMap).Name)
}

func (this *CmdHandler) OnUpdate(obj interface{}, newObj interface{}) {
	fmt.Println("Update: ", newObj.(*v1.ConfigMap).Name)
}

func (this *CmdHandler) OnDelete(obj interface{}) {
	fmt.Println("Delete: ", obj.(*v1.ConfigMap).Name)
}

func main() {
	cliset := config.NewK8sConfig().InitClientset()
	// 通过 clientset 返回一个 listwatcher, 仅支持 default/configmaps 资源
	listWatcher := cache.NewListWatchFromClient(
		cliset.CoreV1().RESTClient(),
		"configmaps",
		"default",
		fields.Everything(),
	)
	// 初始化一个informer, 传入了监听器, 资源名, 间隔同步时间
	// 最后一个是我们定义的 Handler 用于接收我们监听的资源变更事件;
	store, c := cache.NewInformer(listWatcher, &v1.ConfigMap{}, 0, &CmdHandler{})

	// 启动循环监听
	c.Run(wait.NeverStop)

	// 等待3秒 同步缓存
	time.Sleep(3 * time.Second)
	// 从缓存中获取监听到的 configmap 资源
	fmt.Println(store.List())
}
