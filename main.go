package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/kballard/go-shellquote"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kc := kubernetes.NewForConfigOrDie(config)

	// script := "date; echo Hello from the Kubernetes cluster"
	args := []string{
		"date",
		"echo Hello from the Kubernetes cluster",
	}
	script := shellquote.Join(args...)
	fmt.Println(script)

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "arg-demo",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "busybox",
					Image:   "buxybox",
					Command: []string{"/bin/bash"},
					Args: []string{
						"-c",
						script,
					},
					EnvFrom: nil,
					Env: []v1.EnvVar{
						{
							Name:  "A",
							Value: "1",
						},
						{
							Name:  "B",
							Value: "2",
						},
					},
				},
			},
		},
	}
	_, err = kc.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
}
