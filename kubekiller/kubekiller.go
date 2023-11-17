package kubekiller

import (
	"context"
	"fmt"
	"math/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Kubekiller(kubeconfig, namespace string, safemode bool) error {
	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// List the pods in the specified namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		fmt.Println("No pods found in namespace:", namespace)
		return nil
	}

	// Select a random pod
	randPod := pods.Items[rand.Intn(len(pods.Items))]

	if safemode {
		fmt.Printf("Safemode ON: Pod that would have been deleted: %s\n", randPod.Name)
	} else {
		fmt.Printf("Deleting pod: %s\n", randPod.Name)
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), randPod.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
