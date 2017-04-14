// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// Test utilities for kubernetes

const (
	// PodCheckBudget is the maximum number of retries with 1s delays
	PodCheckBudget = 90
)

// CreateNamespace creates a fresh namespace
func CreateNamespace(cl kubernetes.Interface) (string, error) {
	ns, err := cl.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			GenerateName: "istio-test-",
		},
	})
	if err != nil {
		return "", err
	}
	glog.Infof("Created namespace %s", ns.Name)
	return ns.Name, nil
}

// DeleteNamespace removes a namespace
func DeleteNamespace(cl kubernetes.Interface, ns string) {
	if ns != "" && ns != "default" {
		if err := cl.CoreV1().Namespaces().Delete(ns, &meta_v1.DeleteOptions{}); err != nil {
			glog.Warningf("Error deleting namespace: %v", err)
		}
		glog.Infof("Deleted namespace %s", ns)
	}
}

// GetPods gets pod names in a namespace
func GetPods(cl kubernetes.Interface, ns string) []string {
	out := make([]string, 0)
	list, err := cl.CoreV1().Pods(ns).List(meta_v1.ListOptions{})
	if err != nil {
		return out
	}
	for _, pod := range list.Items {
		out = append(out, pod.Name)
	}
	return out
}

// GetAppPods awaits till all pods are running in a namespace, and returns a map
// from "app" label value to the pod name.
func GetAppPods(cl kubernetes.Interface, ns string) (map[string]string, error) {
	pods := make(map[string]string)
	var items []v1.Pod
	for n := 0; ; n++ {
		glog.Info("Checking all pods are running...")
		list, err := cl.CoreV1().Pods(ns).List(meta_v1.ListOptions{})
		if err != nil {
			return pods, err
		}
		items = list.Items
		ready := true

		for _, pod := range items {
			if pod.Status.Phase != "Running" {
				glog.Infof("Pod %s has status %s", pod.Name, pod.Status.Phase)
				ready = false
				break
			}
		}

		if ready {
			break
		}

		if n > PodCheckBudget {
			return pods, fmt.Errorf("exceeded budget for checking pod status")
		}

		time.Sleep(time.Second)
	}

	for _, pod := range items {
		if app, exists := pod.Labels["app"]; exists {
			pods[app] = pod.Name
		}
	}

	return pods, nil
}

// FetchLogs for a container in a a pod
func FetchLogs(cl kubernetes.Interface, name, namespace string, container string) string {
	glog.Infof("Fetching log for container %s in %s.%s", container, name, namespace)
	raw, err := cl.CoreV1().Pods(namespace).
		GetLogs(name, &v1.PodLogOptions{Container: container}).
		Do().Raw()
	if err != nil {
		glog.Infof("Request error %v", err)
		return ""
	}
	return string(raw)
}
