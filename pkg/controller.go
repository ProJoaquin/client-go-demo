package pkg

import (
	"context"
	v17 "k8s.io/api/core/v1"
	v15 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v16 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	netLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"reflect"
	"time"
)

const (
	workNum  = 6
	maxRetry = 6
)

type controller struct {
	client        kubernetes.Interface
	IngressLister netLister.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

func (c controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c controller) updateService(oldobj interface{}, newobj interface{}) {
	// todo 比较Annotation
	if reflect.DeepEqual(oldobj, newobj) {
		return
	}
	c.enqueue(newobj)
}

func (c controller) deleteIngress(obj interface{}) {
	ingress := obj.(*v15.Ingress)
	ownerReference := v16.GetControllerOf(ingress)
	if ownerReference == nil {
		return
	}
	if ownerReference.Kind != "S ervice" {
		return
	}
	c.queue.Add(ingress.Namespace + "/" + ingress.Name)

}

func (c controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceIndexFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	c.queue.Add(key)
}

func (c controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

func (c controller) worker() {
	for c.processNextItem() {

	}
}

func (c controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	key := item.(string)

	err := c.syncService(key)
	if err != nil {
		c.handleError(key, err)
	}
	return true
}

func (c controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	// 是否为资源未找到的错误
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.IngressLister.Ingresses(namespaceKey).Get(name)

	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress
		ig := c.constructIngress(service)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v16.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		err := c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, v16.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c controller) handleError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}

	runtime.HandleError(err)
	c.queue.Forget(key)
}

// 可用yaml嵌入代替construct
func (c controller) constructIngress(service *v17.Service) *v15.Ingress {
	ingress := v15.Ingress{}
	ingress.ObjectMeta.OwnerReferences = []v16.OwnerReference{
		*v16.NewControllerRef(service, v16.SchemeGroupVersion.WithKind("Service")),
	}
	ingress.Name = service.Name
	ingress.Namespace = service.Namespace

	pathType := v15.PathTypePrefix
	icn := "nginx"
	ingress.Spec = v15.IngressSpec{
		IngressClassName: &icn,
		Rules: []v15.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: v15.IngressRuleValue{
					HTTP: &v15.HTTPIngressRuleValue{
						Paths: []v15.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: v15.IngressBackend{
									Service: &v15.IngressServiceBackend{
										Name: service.Name,
										Port: v15.ServiceBackendPort{
											Number: 80,
										},
									},
									Resource: nil,
								},
							},
						},
					},
				},
			},
		},
	}
	return &ingress
}

func Newcontroller(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingressInformer netInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		IngressLister: ingressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}
	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})
	return c
}
