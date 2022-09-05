package webv4

import (
	"fmt"
	"net/http"
	"sync"
)

/*
type Handler interface {
	Routable
	http.Handler
}
*/

/*
// 不使用sync.Map
type HandlerBasedOnMap struct {
	Routable
	handlers map[string]func(c *Context)
}

var _Handler = &HandlerBasedOnMap{}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	// 分发路由

	//	if found {
	//		do
	//	} else {
	//		404
	//	}

	key := h.key(c.R.Method, c.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		// NewContext在Server.Start中创建并加入责任链
		// 不在多路复用器中创建
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
	}
}
*/

// 使用sync.Map
type HandlerBasedOnMap struct {
	Handlers sync.Map
}

var _Handler = &HandlerBasedOnMap{}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	request := c.R
	key := h.key(request.Method, request.URL.Path)
	handler, ok := h.Handlers.Load(key)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		_, _ = c.W.Write([]byte("not any router match"))
		return
	}
	// 类型断言
	handler.(func(c *Context))(c)

}

// 生成key
func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s%s", method, path)
}

func (h *HandlerBasedOnMap) Route(method string, pattern string,
	handerFunc func(c *Context)) {
	key := h.key(method, pattern)

	// 将新的handerFunc加入map中
	//h.handlers[key] = handerFunc
	h.Handlers.Store(key, handerFunc)
}

func NewHanderBasedOnMap() *HandlerBasedOnMap {
	return &HandlerBasedOnMap{}
}
