package webv3

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
	root    Filter
}

func (s *sdkHttpServer) Route(method string, pattern string,
	handlerFunc func(c *Context)) {
	// 将新handerFunc加入map
	s.handler.Route(method, pattern, handlerFunc)
}

/*
func (s *sdkHttpServer) Start(address string) error {

	return http.ListenAndServe(address, s.handler)
}
*/
// 重写Start
func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter,
		request *http.Request) {
		// 不再手写多路复用器
		// NewContext在Start时直接创建并放入责任链
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}

/*
func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name:    name,
		handler: NewHanderBasedOnMap(),
	}
}
*/
// 重写，集成Filter
func NewSdkHttpServer(name string, builders ...FilterBuilder) Server {

	handler := NewHanderBasedOnMap()

	/*
		// 修改Hander接口前
			var root Filter = func(c *Context) {
				handler.ServeHTTP(c.W, c.R)
			}
	*/

	//
	var root Filter = handler.ServeHTTP

	// 连接责任链
	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	res := &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
	return res
}
