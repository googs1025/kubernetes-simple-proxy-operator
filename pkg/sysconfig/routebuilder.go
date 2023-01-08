package sysconfig

import (
	"github.com/gorilla/mux"
	"net/http"
)

// route构建器， build 方法必须要执行
var MyRouter *mux.Router

func init() {
	MyRouter = mux.NewRouter()
}


type RouteBuilder struct {
	route *mux.Route
}

func NewRouteBuilder() *RouteBuilder {
	return &RouteBuilder{route: MyRouter.NewRoute()}
}

func(rb *RouteBuilder) SetPath(path string ,exact bool) *RouteBuilder{

	if exact {
		rb.route.Path(path)
	} else {
		rb.route.PathPrefix(path)
	}
   return rb
}

// 第二个参数是故意的，方便调用时 传入 条件，省的外面写 if else
func(rb *RouteBuilder) SetHost(host string, set bool) *RouteBuilder{
	if set {
		rb.route.Host(host)
	}
	 return rb
}

func(rb *RouteBuilder) Build(handler http.Handler)  {
	rb.route.
		Methods("GET","POST","PUT","DELETE","OPTIONS").
		Handler(handler)
}
