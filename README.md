# **Kubernetes Informers Examples**
This repository contains examples of using Kubernetes Informers with Go. Informers are a key component of the Kubernetes client-go library, allowing efficient watching and caching of Kubernetes resources.

# Examples
## 1. Pod Informer
The Pod Informer example demonstrates how to watch for changes to Pods in a Kubernetes cluster.
Key features:

Watches for Pod creation, updates, and deletions
Logs Pod events with namespace and name information
Uses the coreinformers.PodInformer

Use this example when you need to monitor or react to changes in Pod resources.
## 2. Deployment Informer
The Deployment Informer example shows how to watch for changes to Deployments in a Kubernetes cluster.
Key features:

Watches for Deployment creation, updates, and deletions
Logs Deployment events with namespace, name, and replica information
Uses the appsinformers.DeploymentInformer

This example is useful when you need to track changes in Deployment resources, such as scaling events or updates to the Deployment specification.
## 3. Multi-resource Informer
The Multi-resource Informer example demonstrates how to watch for changes to multiple resource types (Deployments, Services, and Pods) in a single controller.

## Key features:
- Watches for changes in Deployments, Services, and Pods
- Logs events for all three resource types
- Uses multiple informers: appsinformers.DeploymentInformer, coreinformers.ServiceInformer, and coreinformers.PodInformer
- Demonstrates how to combine multiple informers in a single controller

This example is ideal when you need to monitor or react to changes across different types of Kubernetes resources in a coordinated manner.

## Usage
For each example:

- Ensure you have the necessary Kubernetes Go client libraries installed.
- Set up your Kubernetes configuration file (default location: ~/.kube/config).
- Run the program using go run main.go.

The program will start watching for changes to the specified resources in your Kubernetes cluster and log events as they occur.
