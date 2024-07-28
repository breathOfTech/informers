package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type MultiResourceController struct {
	informerFactory informers.SharedInformerFactory
	deployInformer  appsinformers.DeploymentInformer
	serviceInformer coreinformers.ServiceInformer
	podInformer     coreinformers.PodInformer
}

func (c *MultiResourceController) Run(stopCh <-chan struct{}) error {
	c.informerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh,
		c.deployInformer.Informer().HasSynced,
		c.serviceInformer.Informer().HasSynced,
		c.podInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}

func (c *MultiResourceController) deploymentAdd(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	klog.Infof("DEPLOYMENT CREATED: %s/%s", deploy.Namespace, deploy.Name)
}

func (c *MultiResourceController) deploymentUpdate(old, new interface{}) {
	oldDeploy := old.(*appsv1.Deployment)
	newDeploy := new.(*appsv1.Deployment)
	klog.Infof("DEPLOYMENT UPDATED: %s/%s, Replicas: %d", oldDeploy.Namespace, oldDeploy.Name, newDeploy.Status.Replicas)
}

func (c *MultiResourceController) deploymentDelete(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	klog.Infof("DEPLOYMENT DELETED: %s/%s", deploy.Namespace, deploy.Name)
}

func (c *MultiResourceController) serviceAdd(obj interface{}) {
	svc := obj.(*corev1.Service)
	klog.Infof("SERVICE CREATED: %s/%s", svc.Namespace, svc.Name)
}

func (c *MultiResourceController) serviceUpdate(old, new interface{}) {
	oldSvc := old.(*corev1.Service)
	newSvc := new.(*corev1.Service)
	klog.Infof("SERVICE UPDATED: %s/%s, Type: %s", oldSvc.Namespace, oldSvc.Name, newSvc.Spec.Type)
}

func (c *MultiResourceController) serviceDelete(obj interface{}) {
	svc := obj.(*corev1.Service)
	klog.Infof("SERVICE DELETED: %s/%s", svc.Namespace, svc.Name)
}

func (c *MultiResourceController) podAdd(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("POD CREATED: %s/%s", pod.Namespace, pod.Name)
}

func (c *MultiResourceController) podUpdate(old, new interface{}) {
	oldPod := old.(*corev1.Pod)
	newPod := new.(*corev1.Pod)
	klog.Infof("POD UPDATED: %s/%s, Status: %s", oldPod.Namespace, oldPod.Name, newPod.Status.Phase)
}

func (c *MultiResourceController) podDelete(obj interface{}) {
	pod := obj.(*corev1.Pod)
	klog.Infof("POD DELETED: %s/%s", pod.Namespace, pod.Name)
}

func NewMultiResourceController(informerFactory informers.SharedInformerFactory) (*MultiResourceController, error) {
	deployInformer := informerFactory.Apps().V1().Deployments()
	serviceInformer := informerFactory.Core().V1().Services()
	podInformer := informerFactory.Core().V1().Pods()

	c := &MultiResourceController{
		informerFactory: informerFactory,
		deployInformer:  deployInformer,
		serviceInformer: serviceInformer,
		podInformer:     podInformer,
	}

	_, err := deployInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.deploymentAdd,
			UpdateFunc: c.deploymentUpdate,
			DeleteFunc: c.deploymentDelete,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = serviceInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.serviceAdd,
			UpdateFunc: c.serviceUpdate,
			DeleteFunc: c.serviceDelete,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.podAdd,
			UpdateFunc: c.podUpdate,
			DeleteFunc: c.podDelete,
		},
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "absolute path to the kubeconfig file")
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error creating clientset: %s", err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller, err := NewMultiResourceController(factory)
	if err != nil {
		klog.Fatalf("Error creating controller: %s", err.Error())
	}

	stop := make(chan struct{})
	defer close(stop)
	err = controller.Run(stop)
	if err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
	select {}
}
