package webv3

type Handler interface {
	Routable
	ServeHTTP(c *Context)
}
