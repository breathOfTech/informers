package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	dirname, _ := os.UserHomeDir()
	kubeConfig := flag.String("kubeconfig", filepath.Join(dirname, ".kube/config"), "location of your config file")
	// Build the Kubernetes client config
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		// Handle error
		fmt.Printf("error %s,  building config from flags", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting incluster config", err.Error())
			panic(err.Error())
		}
	}

	// Create the kubernetes client
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Create a shared informer
	informerFactory := informers.NewSharedInformerFactory(clientSet, 30*time.Second)
	informers.NewFilteredSharedInformerFactory(clientSet, 30*time.Second, "my-namespace", func(lo *v1.ListOptions) {
		lo.LabelSelector = ""
	})

	podInformer := informerFactory.Core().V1().Pods().Informer()

	// Event handler
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("New Pod Added: %s\n", pod.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPod := oldObj.(*corev1.Pod)
			newPod := newObj.(*corev1.Pod)
			fmt.Printf("Pod Updated: %s -> %s\n", oldPod.Name, newPod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("Pod Deleted: %s\n", pod.Name)
		},
	})

	// Start and WaitForCacheSync
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	// This call goes to cache and not the K8s API server
	pod, _ := informerFactory.Core().V1().Pods().Lister().Pods("default").Get("default")
	fmt.Println(pod)
}
