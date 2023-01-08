package sysconfig

import (
	"io/ioutil"
	"k8s.io/api/networking/v1"
	"log"
	"sigs.k8s.io/yaml"
)

type Server struct {
	Port int   //代表是代理启动端口
}

// app.yaml的配置文件内容
type SysConfigStruct struct {
	Server Server
	Ingress []v1.Ingress
}



func NewSysConfig() *SysConfigStruct {
	return &SysConfigStruct{
	}
}

var SysConfig *SysConfigStruct

func InitConfig()  {
	// 读取yaml配置
	config, err := ioutil.ReadFile("./app.yaml")
	if err != nil {
		log.Fatal(err)
	}

	SysConfig = NewSysConfig()

	err = yaml.Unmarshal(config, SysConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 解析配置文件
    ParseRule()

}
