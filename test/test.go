package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/yeqown/log"
	"k8s.io/klog/v2"
	"try_operator/pkg/filters"
	"try_operator/pkg/sysconfig"
)


func ProxyHandler(ctx *fasthttp.RequestCtx){
	// 代表匹配到了 path
	 if getProxy := sysconfig.GetRoute(ctx.Request); getProxy != nil {
		 filters.ProxyFilters(getProxy.RequestFilters).Do(ctx) //过滤
		 getProxy.Proxy.ServeHTTP(ctx) // 反代
		 filters.ProxyFilters(getProxy.ResponseFilters).Do(ctx) //过滤
	 } else { // 404代表没有匹配到
		 ctx.Response.SetStatusCode(404)
		 ctx.Response.SetBodyString("404...")
	 }

}

func main() {
	// 初始化配置文件
	err := sysconfig.InitConfig()
	if err != nil {
		klog.Error("config error")
		return
	}
	klog.Info("start server!")
	// 启动http
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", sysconfig.SysConfig.Server.Port), ProxyHandler))
}
