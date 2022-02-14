package main

import (
	"context"
	"fmt"
	"log"
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

	zero := int64(0)

	clientset, err := kubernetes.NewForConfig(k8sConfig)
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	for _, p := range pods.Items {
		gracePeriodSeconds := time.Duration(30)
		if p.Spec.TerminationGracePeriodSeconds != nil {
			gracePeriodSeconds = time.Duration(*p.Spec.TerminationGracePeriodSeconds)
		}
		if p.ObjectMeta.DeletionTimestamp != nil && time.Now().Sub(p.ObjectMeta.DeletionTimestamp.Time) > ((gracePeriodSeconds*time.Second)+(5*time.Minute)) {
			fmt.Printf("Killing: %s/%s\n", p.Namespace, p.Name)
			err = clientset.CoreV1().Pods(p.Namespace).Delete(context.Background(), p.Name, metav1.DeleteOptions{GracePeriodSeconds: &zero})
			if err != nil {
				fmt.Printf("Got error when kill %s/%s. Err: %s\n", p.Namespace, p.Name, err.Error())
			}
		}
	}
}
