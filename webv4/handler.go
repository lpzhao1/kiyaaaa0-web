package webv4

type Handler interface {
	Routable
	ServeHTTP(c *Context)
}
