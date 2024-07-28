package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/informers"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/tools/cache"
)

// ServiceLoggingController logs the name and namespace of services that are added,
// deleted, or updated
type DeploymentLoggingController struct {
	informerFactory    informers.SharedInformerFactory
	deploymentInformer appsinformers.DeploymentInformer
}

func (c *DeploymentLoggingController) deploymentAdd(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	fmt.Printf("DEPLOYMENT CREATED: %s/%s", deploy.Namespace, deploy.Name)
}

func (c *DeploymentLoggingController) deploymentUpdate(old, new interface{}) {
	oldDeploy := old.(*appsv1.Deployment)
	newDeploy := new.(*appsv1.Deployment)
	fmt.Printf("DEPLOYMENT UPDATED: %s/%s, Replicas: %d", oldDeploy.Namespace, oldDeploy.Name, newDeploy.Status.Replicas)
}

func (c *DeploymentLoggingController) deploymentDelete(obj interface{}) {
	deploy := obj.(*appsv1.Deployment)
	fmt.Printf("DEPLOYMENT DELETED: %s/%s", deploy.Namespace, deploy.Name)
}

// NewServiceLoggingController creates a ServiceLoggingController
func NewDeploymentLoggingController(informerFactory informers.SharedInformerFactory) (*DeploymentLoggingController, error) {
	deploymentInformer := informerFactory.Apps().V1().Deployments()

	c := &DeploymentLoggingController{
		informerFactory:    informerFactory,
		deploymentInformer: deploymentInformer,
	}
	_, err := deploymentInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.deploymentAdd,
			// Called on resource update and every resyncPeriod on existing resources.
			UpdateFunc: c.deploymentUpdate,
			// Called on resource deletion.
			DeleteFunc: c.deploymentDelete,
		},
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *DeploymentLoggingController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so
	// far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.deploymentInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}
