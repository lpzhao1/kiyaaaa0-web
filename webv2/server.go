package webv2

import "net/http"

/*
type sdkHttpServer struct {
	Name string
}

// 原始版本
func (s *sdkHttpServer) Route(pattern string, handlefunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlefunc)
}



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
*/

type Routable interface {
	Route(method string, pattern string, handerFunc func(*Context))
}

type Server interface {
	Routable
	Start(address string) error
}

type sdkHttpServer struct {
	Name string
	// 自己实现一个Handler负责路由
	handler Handler
}

func (s *sdkHttpServer) Route(method string, pattern string,
	handlerFunc func(c *Context)) {
	// 将新handerFunc加入map
	s.handler.Route(method, pattern, handlerFunc)
}

//
func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, s.handler)
}

func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name:    name,
		handler: NewHanderBasedOnMap(),
	}
}
