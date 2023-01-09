package k8sconfig

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"try_operator/pkg/common"
)

func K8sRestConfig() *rest.Config {
	// 读取配置
	path := common.GetWd()
	config, err := clientcmd.BuildConfigFromFlags("", path + "/resource/config")
	if err != nil {
		log.Fatal(err)
	}
	config.Insecure = true

	return config

}
