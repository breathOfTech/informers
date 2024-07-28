package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "absolute path to the kubeconfig file")
}
func main() {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	controller, err := NewDeploymentLoggingController(factory)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan struct{})
	defer close(stop)
	err = controller.Run(stop)
	if err != nil {
		log.Fatal(err)
	}
	select {}
}
