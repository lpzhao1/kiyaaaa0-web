package webv5

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

/*
// 原始版本
func (s *sdkHttpServer) Route(pattern string, handlefunc http.HandlerFunc) {
	http.HandleFunc(pattern, handlefunc)
}

//
func (s *sdkHttpServer) Start(address string) error {
	return http.ListenAndServe(address, nil)
}
*/

/*
// 为了让路由支持由 web 框架创建 context 进行改造
func (s *sdkHttpServer) Route(pattern string, handlerfunc func(c *Context)) {
	http.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		// 先创建 context 再执行 handlerfunc
		c := NewContext(writer, request)
		handlerfunc(c)
	})
}

func (s *sdkHttpServer) Start(address string) error {
	// s.handler 作为手写处理器
	return http.ListenAndServe(address, s.handler)
}

// 加入责任链重写Start
func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter,
		request *http.Request) {
		// 不再手写处理器
		// NewContext在Start时直接创建并放入责任链
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}
*/

// Server 和 Handler 都实现了Route，进一步抽象
type Routable interface {
	Route(method string, pattern string, handerFunc func(*Context))
}

// 对 Server 的顶级抽象
type Server interface {
	// 路由
	Routable
	// 开启服务
	Start(address string) error
	// v4添加 Shutdown
	Shutdown(ctx context.Context) error
}

type sdkHttpServer struct {
	Name string
	// 自己实现一个 Handler 负责路由
	handler Handler
	root    Filter
	// 添加 sync.Pool
	ctxPool sync.Pool
}

func (s *sdkHttpServer) Route(method string, pattern string,
	handlerFunc func(c *Context)) {
	// Server 与 Handler 都实现了Route
	// 直接调用 Handler 的 Route
	s.handler.Route(method, pattern, handlerFunc)
}

// 再重写Start
func (s *sdkHttpServer) Start(address string) error {
	/*
		该版本中 newContext 不在Start 函数中创建
		直接在 NewSdkHttpServer 时初始化
	*/
	// 将 sdkHttpServer 作为自写处理器
	/*
		http.Handler 是一个 interface{}
		s 实现了 http.Handler 的 ServeHTTP方法
		可以作为 http.ListenAndServe 的参数
	*/
	return http.ListenAndServe(address, s)
}

func (s *sdkHttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 为什么要这么写
	/*
		ctxPool 起始时为空，Get 方法得到一个 interface{}
		类型断言后即为 newContext 返回的指向空 Context 的指针
		Reset 方法将其字段初始化
		由责任链根部开始执行
		最后将其放回 ctxPool
	*/
	c := s.ctxPool.Get().(*Context)
	defer func() {
		s.ctxPool.Put(c)
	}()
	c.Reset(writer, request)
	s.root(c)
}

func (s *sdkHttpServer) Shutdown(ctx context.Context) error {
	// 因为我们这个简单的框架，没有什么要清理的，
	// 所以我们 sleep 一下来模拟这个过程
	fmt.Printf("%s shutdown...\n", s.Name)
	time.Sleep(time.Second)
	fmt.Printf("%s shutdown!!!\n", s.Name)
	return nil
}

/*
// 加入责任链前
func NewSdkHttpServer(name string) Server {
	return &sdkHttpServer{
		Name:    name,
		handler: NewHanderBasedOnMap(),
	}
}
*/
// 加入责任链，重写，集成 Filter
func NewSdkHttpServer(name string, builders ...FilterBuilder) Server {

	// handler := NewHanderBasedOnMap()
	// 改用路由树了
	handler := NewHanderBasedOnTree()

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
		ctxPool: sync.Pool{New: func() interface{} {
			// newContext不在Start函数中创建
			// 直接在NewSdkHttpServer时初始化
			return newContext()
		}},
	}
	return res
}
