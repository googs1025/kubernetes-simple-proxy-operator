package sysconfig

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	networkingv1 "k8s.io/api/networking/v1"
	"os"
	"sigs.k8s.io/yaml"
	"try_operator/pkg/common"
)

type Server struct {
	Port int   //代表是代理启动端口
}

// app.yaml的配置文件内容
type SysConfigStruct struct {
	Server Server
	Ingress []networkingv1.Ingress
}



func NewSysConfig() *SysConfigStruct {
	return &SysConfigStruct{
	}
}

var SysConfig *SysConfigStruct

func InitConfig() error {
	// 读取yaml配置
	config, err := ioutil.ReadFile("./app.yaml")
	if err != nil {
		return err
	}

	SysConfig = NewSysConfig()

	err = yaml.Unmarshal(config, SysConfig)
	if err != nil {
		return err
	}

	// 解析配置文件
    ParseRule()

	return nil

}

func AppConfig(ingress *networkingv1.Ingress) error {
	isEdit := false

	// 更新内存的配置
	for i, config := range SysConfig.Ingress {
		// 能在内存找到，代表是更新
		if config.Name == ingress.Name && config.Namespace == ingress.Namespace {
			SysConfig.Ingress[i] = *ingress
			isEdit = true
			break
		}
	}

	// 新加入的
	if !isEdit {
		SysConfig.Ingress = append(SysConfig.Ingress, *ingress)

	}

	if err := saveConfigToFile(); err != nil {
		return err
	}


	return ReloadConfig()

}

// ReloadConfig 重载配置
func ReloadConfig() error {
	MyRouter = mux.NewRouter()
	return InitConfig()

}

func DeleteConfig(name, namespace string) error {
	isEdit := false
	for i, config := range SysConfig.Ingress {
		if config.Name == name && config.Namespace == namespace {
			SysConfig.Ingress = append(SysConfig.Ingress[:i], )
			isEdit = true
			break
		}
	}
	if isEdit {
		if err := saveConfigToFile(); err != nil {
			return err
		}
		return ReloadConfig()
	}

	return nil
}

func saveConfigToFile() error {

	b, err := yaml.Marshal(SysConfig)
	if err != nil {
		return err
	}
	// 读取文件
	path := common.GetWd()
	filePath := path + "/app.yaml"
	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
	if err != nil {
		return err
	}

	defer appYamlFile.Close()
	_, err = appYamlFile.Write(b)
	if err != nil {
		return err
	}

	return nil
}