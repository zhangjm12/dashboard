// Copyright 2015 Google Inc. All Rights Reserved.
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

package client

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/client-go/kubernetes"
	api "k8s.io/client-go/pkg/api/v1"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/kube"
)

const (
	tillerNamespace = "kube-system"
	tillerPort      = 44134
)

func CreateHelmTillerClient(apiclient *kubernetes.Clientset) (*helm.Client, error) {
	tunnel, err := newTillerPortForwarder(tillerNamespace, apiclient)
	if err != nil {
		return nil, err
	}
	log.Printf("Created tunnel using local port: '%d'", tunnel.Local)
	tillerHost := fmt.Sprintf(":%d", tunnel.Local)
	log.Printf("Creating tiller client using host: %q", tillerHost)
	tillerClient := helm.NewClient(helm.Host(tillerHost))
	return tillerClient, nil
}

// TODO: refactor out this global var
var tunnel *kube.Tunnel

func newTillerPortForwarder(namespace string, apiclient *kubernetes.Clientset) (*kube.Tunnel, error) {
	podName, err := getTillerPodName(apiclient, namespace)
	if err != nil {
		return nil, err
	}
	log.Printf("tiller pod found: %q", podName)

	config, err := kube.GetConfig("").ClientConfig()
	if err != nil {
		return nil, err
	}
	t := kube.NewTunnel(apiclient.CoreV1().RESTClient(), config, namespace, podName, tillerPort)
	return t, t.ForwardPort()
}

func getTillerPodName(client *kubernetes.Clientset, namespace string) (string, error) {
	// TODO: use a const for labels
	selector := labels.Set{"app": "helm", "name": "tiller"}.AsSelector()
	pod, err := getFirstRunningPod(client, namespace, selector)
	if err != nil {
		return "", err
	}
	return pod.ObjectMeta.GetName(), nil
}

func getFirstRunningPod(client *kubernetes.Clientset, namespace string, selector labels.Selector) (*api.Pod, error) {
	options := metav1.ListOptions{LabelSelector: selector.String()}
	pods, err := client.CoreV1().Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, fmt.Errorf("could not find tiller")
	}
	for _, p := range pods.Items {
		if api.IsPodReady(&p) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find a ready pod")
}
