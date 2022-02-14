package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"log"
	"os"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Get in cluster k8sConfig error: %s\n\n", err.Error())
		log.Println("Will try connect to local 8001 (kubectl proxy)")
		k8sConfig = &rest.Config{
			Host: "http://127.0.0.1:8001",
		}
	}

	// Namespace to check. If empty, it check all namespace. Can use ',' to separate namespaces.
	targetNamespace := os.Getenv("NAMESPACE")

	zero := int64(0)

	clientset, err := kubernetes.NewForConfig(k8sConfig)

	var pods []v1.Pod

	if strings.TrimSpace(targetNamespace) == "" {
		podList, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
		pods = append(pods, podList.Items...)
	} else {
		for _, ns := range strings.Split(targetNamespace, ",") {
			ns := strings.TrimSpace(ns)
			if ns != "" {
				podList, _ := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
				pods = append(pods, podList.Items...)
			}
		}
	}
	for _, p := range pods {
		gracePeriodSeconds := time.Duration(30)
		if p.Spec.TerminationGracePeriodSeconds != nil {
			gracePeriodSeconds = time.Duration(*p.Spec.TerminationGracePeriodSeconds)
		}
		if p.ObjectMeta.DeletionTimestamp != nil && time.Now().Sub(p.ObjectMeta.DeletionTimestamp.Time) > ((gracePeriodSeconds*time.Second)+(5*time.Minute)) {
			log.Printf("Killing: %s/%s\n", p.Namespace, p.Name)
			err = clientset.CoreV1().Pods(p.Namespace).Delete(context.Background(), p.Name, metav1.DeleteOptions{GracePeriodSeconds: &zero})
			if err != nil {
				log.Printf("Got error when kill %s/%s. Err: %s\n", p.Namespace, p.Name, err.Error())
			}
		}
	}
}
