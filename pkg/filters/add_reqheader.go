package filters

import (
	"github.com/valyala/fasthttp"
	"strings"
)
const AddRequestHeaderAnnotation = AnnotationPrefix + "/add-request-header"

func init() {
	registerFilter(AddRequestHeaderAnnotation,(*AddRequestHeaderFilter)(nil))
}

type AddRequestHeaderFilter struct {
    pathValue string
    target string  //注解  值
    path string
}


func(a *AddRequestHeaderFilter) SetPath(value  string){}

// 可变参数。第1个是注解值:的值 如 /$1
func(a *AddRequestHeaderFilter) SetValue(values ...string){
	a.target = values[0]
}

func(a *AddRequestHeaderFilter) Do(ctx *fasthttp.RequestCtx){
	 kvList := strings.Split(a.target,";")
	 for _, kv := range kvList {
	 	keyAndValue := strings.Split(kv,"=")
	 	if len(keyAndValue) == 2 {
			ctx.Request.Header.Add(keyAndValue[0],keyAndValue[1])
		}
	 }

}