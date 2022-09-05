package webv4

type HandlerBasedOnTree struct {
	root *node
}

type node struct {
	path     string
	children []*node

	// 如果是叶子节点
	// 匹配上后就可以调用该方法
	// handler HandlerFunc
	handler func(c *Context)
}
