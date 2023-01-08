package k8sconfig

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func K8sRestConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/zhenyu.jiang/go/src/golanglearning/new_project/try_operator/resource/config")
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true

	return config

}
