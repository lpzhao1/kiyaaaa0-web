package webv5

type Handler interface {
	Routable
	ServeHTTP(c *Context)
}
