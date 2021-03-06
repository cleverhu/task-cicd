package k8sconfig

import (
	"log"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//全局变量

const NSFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

//POD里  体内
func K8sRestConfigInPod() *rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	return config
}

// 获取 config对象
func K8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "./resources/config")
	if err != nil {
		return K8sRestConfigInPod()
	}
	return config
}
