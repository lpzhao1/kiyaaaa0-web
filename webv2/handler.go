package webv2

import (
	"fmt"
	"net/http"
)

type Handler interface {
	Routable
	http.Handler
}

type HandlerBasedOnMap struct {
	Routable
	handlers map[string]func(c *Context)
}

var _Handler = &HandlerBasedOnMap{}

func (h *HandlerBasedOnMap) ServeHTTP(writer http.ResponseWriter,
	request *http.Request) {
	// 分发路由
	/*
		if found {
			do
		} else {
			404
		}
	*/
	key := h.key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {

		// NewContext不再Server.Route时创建
		// 在Server.Start调用自己重写的多路复用器时
		// 确认handler存在再创建
		c := NewContext(writer, request)
		handler(c)
	} else {
		writer.WriteHeader(http.StatusNotFound)
		_, _ = writer.Write([]byte("not any router match"))
	}
}

// 生成key
func (h *HandlerBasedOnMap) key(method string, path string) string {
	return fmt.Sprintf("%s%s", method, path)
}

func (h *HandlerBasedOnMap) Route(method string, pattern string,
	handerFunc func(c *Context)) {
	key := h.key(method, pattern)
	// 将新的handerFunc加入map中
	h.handlers[key] = handerFunc
}

func NewHanderBasedOnMap() Handler {
	return &HandlerBasedOnMap{
		handlers: make(map[string]func(c *Context), 32),
	}
}
