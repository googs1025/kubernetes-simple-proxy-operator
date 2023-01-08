package filters

import (
	"github.com/valyala/fasthttp"
	"strings"
)
const AddResponseHeaderAnnotation = AnnotationPrefix + "/add-response-header"

func init() {
	//注册响应 过滤器
	registerResponseFilter(AddResponseHeaderAnnotation,(*AddResponseHeaderFilter)(nil))
}

type AddResponseHeaderFilter struct {
    pathValue string
    target string  //注解  值
    path string
}

func(a *AddResponseHeaderFilter) SetPath(value  string){}

// 可变参数。第1个是 注解值:的值 如 /$1
func(a *AddResponseHeaderFilter) SetValue(values ...string){
	a.target=values[0]
}

func(a *AddResponseHeaderFilter) Do(ctx *fasthttp.RequestCtx){
	 kvList := strings.Split(a.target,";")
	 for _, kv := range kvList {
	 	keyAndValue := strings.Split(kv,"=")
	 	if len(keyAndValue) == 2 {
			ctx.Response.Header.Add(keyAndValue[0], keyAndValue[1])
		}
	 }

}