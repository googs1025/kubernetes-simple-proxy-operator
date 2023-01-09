package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"try_operator/pkg/controller"
	"try_operator/pkg/filters"
	"try_operator/pkg/k8sconfig"
	"try_operator/pkg/sysconfig"
)

/*
	manager 主要用来管理Controller Admission Webhook 包括：
	访问资源对象的client cache scheme 并提供依赖注入机制 优雅关闭机制

	operator = crd + controller + webhook
*/

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

	mgr, err := manager.New(k8sconfig.K8sRestConfig(), manager.Options{})
	if err  != nil {
		klog.Error(err, "unable to set up manager")
		os.Exit(1)
	}

	proxyCtl := controller.NewProxyController()
	// 传入资源&v1.Ingress{}，也可以用crd
	err = builder.ControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Watches(&source.Kind{ // 加入监听。
			Type: &networkingv1.Ingress{},
		}, handler.Funcs{
			DeleteFunc: proxyCtl.IngressDeleteHandler,
		}).
		Complete(controller.NewProxyController())

	//++ 注册进入序列化表
	err = k8sconfig.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		klog.Error(err, "unable add schema")
		os.Exit(1)
	}

	// 载入业务配置
	if err = sysconfig.InitConfig(); err != nil {
		klog.Error(err, "unable to load sysconfig")
		os.Exit(1)
	}
	errC := make(chan error)

	// 启动controller管理器
	go func() {
		klog.Info("controller start!! ")
		if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
			errC <-err
		}
	}()

	// 启动网关
	go func() {
		klog.Info("proxy start!! ")
		if err = fasthttp.ListenAndServe(fmt.Sprintf(":%d", sysconfig.SysConfig.Server.Port), ProxyHandler);err!=nil{
			errC <-err
		}
	}()

	getError := <-errC
	log.Println(getError.Error())

}


