package kubekiller

import "fmt"

func Kubekiller(kubeconfig, namespace string) error {
	fmt.Print("Hello ", kubeconfig, "\n")
	fmt.Print("Hello ", namespace, "\n")
	return nil
}
