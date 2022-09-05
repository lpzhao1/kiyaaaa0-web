package webv1

import (
	"net/http"
)

type Server interface {
	Route(pattern string, handlefunc func(*Context))
	Start(address string) error
}

/*
func (s *sdkHttpServer) Route(pattern string, handlefunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlefunc)
}
*/

// 为了让路由支持由web框架创建context进行改造
func (s *sdkHttpServer) Route(pattern string, handlerfunc func(c *Context)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		//先创建context再执行handlerfunc
		c := NewContext(writer, request)
		handlerfunc(c)
	})
}

func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, nil)
}

type sdkHttpServer struct {
	Name string
}

func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}
