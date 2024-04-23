package main

import (
	"fmt"
	"k8s.io/client-go/tools/cache"
)

type Pod struct {
	Name  string
	Value int
}

func NewPod(name string, v int) Pod {
	return Pod{Name: name, Value: v}
}

// 需要提供一个资源的唯一标识的字符串给到 DeltaFifo， 这样它就能追踪某个资源的变化
func PodKeyFunc(obj interface{}) (string, error) {
	return obj.(Pod).Name, nil
}

func main() {
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{KeyFunction: PodKeyFunc})

	// ADD3个object 进入 fifo
	pod1 := NewPod("pod-1", 1)
	pod2 := NewPod("pod-2", 2)
	pod3 := NewPod("pod-3", 3)
	df.Add(pod1)
	df.Add(pod2)
	df.Add(pod3)
	// Update pod-1
	pod1.Value = 11
	df.Update(pod1)
	df.Delete(pod1)

	// 当前df 的列表
	fmt.Println(df.List())

	// 循环抛出事件
	for {
		df.Pop(func(i interface{}) error {
			for _, delta := range i.(cache.Deltas) {
				switch delta.Type {
				case cache.Added:
					fmt.Printf("Add Event: %v \n", delta.Object)
					break
				case cache.Updated:
					fmt.Printf("Update Event: %v \n", delta.Object)
					break
				case cache.Deleted:
					fmt.Printf("Delete Event: %v \n", delta.Object)
					break
				case cache.Sync:
					fmt.Printf("Sync Event: %v \n", delta.Object)
					break
				case cache.Replaced:
					fmt.Printf("Replaced Event: %v \n", delta.Object)
					break
				}
			}
			return nil
		})
	}

}
