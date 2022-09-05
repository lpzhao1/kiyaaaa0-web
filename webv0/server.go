package webv0

import (
	"net/http"
)

type Server interface {
	Route(pattern string, handlefunc http.HandlerFunc)
	Start(address string) error
}

func (s *sdkHttpServer) Route(pattern string, handlefunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlefunc)
}

/*
// 为了让路由支持由web框架创建context进行改造
func (s *sdkHttpServer) Route(pattern string, handlefunc func(*Context)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter,
		request *http.Request) {
		c := NewContext(writer, request)
		HandleFunc(c)
	})
}
*/

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
