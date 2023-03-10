package sysconfig

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/valyala/fasthttp"
	"github.com/yeqown/fasthttp-reverse-proxy/v2"
	"try_operator/pkg/filters"
	v1 "k8s.io/api/networking/v1"
	"net/http"
	"net/url"
)
type ProxyHandler struct {
	Proxy *proxy.ReverseProxy  // proxy对象。 保存proxy
	RequestFilters []filters.ProxyFilter
	ResponseFilters []filters.ProxyFilter
}

// 空函数没用
func(p *ProxyHandler) ServeHTTP(http.ResponseWriter, *http.Request){}

//解析配置文件中的rules， 初始化 路由
func ParseRule()  {
	//三次循环遍历，解析切片
	for _ , ingress := range SysConfig.Ingress {
		for _ , rule := range ingress.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				//构建 ReverseProxy代理对象
				rProxy := proxy.NewReverseProxy(
					fmt.Sprintf("%s:%d", path.Backend.Service.Name, path.Backend.Service.Port.Number))
				// 建造者模式使用
				routeBud := NewRouteBuilder()

			    routeBud.
					SetPath(path.Path,path.PathType != nil && *path.PathType == v1.PathTypeExact).
					SetHost(rule.Host, rule.Host != "").
					Build(&ProxyHandler{Proxy: rProxy,
						RequestFilters: filters.CheckAnnotations(ingress.Annotations,false),
						ResponseFilters: filters.CheckAnnotations(ingress.Annotations,true),
				})

			}
		}
	}

}

// 获取路由（先匹配 请求path，如果匹配到，会返回 对应的proxy对象)
func GetRoute(req fasthttp.Request)*ProxyHandler{
	match := &mux.RouteMatch{}
	// 请求需要用的request
	httpReq := &http.Request{
		URL: &url.URL{Path:  string(req.URI().Path())},
		Method: string(req.Header.Method()),
		Host: string(req.Header.Host()),
	}

	if MyRouter.Match(httpReq,match){  // 匹配到

		proxyHandler := match.Handler.(*ProxyHandler)
		pathExp,err := match.Route.GetPathRegexp()  //对过滤器 放入值：path
		// 譬如这样：^/users/(?P<v0>[^/]+)
		if err == nil {
			//不管是 Request还是Response都要设置path
			filters.ProxyFilters(proxyHandler.RequestFilters).SetPath(pathExp)
			filters.ProxyFilters(proxyHandler.ResponseFilters).SetPath(pathExp)
		}
		return proxyHandler

	}
	return  nil
}
