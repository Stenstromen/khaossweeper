package kubekiller

import (
	"context"
	"fmt"
	"math/rand"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Kubekiller(kubeconfig, namespace string, safemode bool) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", err
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no pods found in namespace: %s", namespace)
	}

	randPod := pods.Items[rand.Intn(len(pods.Items))]

	if safemode {
		return fmt.Sprintf("(Would have) %s", randPod.Name), nil
	}

	fmt.Printf("Deleting pod: %s\n", randPod.Name)
	err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), randPod.Name, metav1.DeleteOptions{})
	if err != nil {
		return randPod.Name, err
	}

	return randPod.Name, nil
}
