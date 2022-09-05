package webv5

import "strings"

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

func NewHanderBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: &node{},
	}
}

// ServeHTTP 就是从树中找节点
// 找到了就执行
func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	// panic("implement me")
	handler, found := h.findRouter(c.R.URL.Path)
	if !found {
		c.W.WriteHeader(404)
		_, _ = c.W.Write([]byte("Not Found"))
		return
	}
	handler(c)

	// old version
	/*
		url := strings.Trim(c.R.URL.Path, "/")
		paths := strings.Split(url, "/")
		cur := h.root
		for _, path := range paths {
			// 从子节点中找一个匹配到当前 path 的节点
			mathChild, found := h.findMatchChild(cur, path)
			if !found {
				// 找不到匹配的路径，直接返回
				c.W.WriteHeader(404)
				_, _ = c.W.Write([]byte("Not Found"))
				return
			}
			cur = mathChild
		}

		// 这里本来应该已经找完了
		if cue.handler == nil {
			// 这里适用以下场景
			// 已经注册了 /user/friends
			// 然后访问 /user
			c.W.WriteHeader(404)
			_, _ = c.W.Write([]byte("Not Found"))
			return
		}

		//
		cur.handler(c)
	*/
}

func (h *HandlerBasedOnTree) findRouter(path string) (func(c *Context), bool) {
	// 去除头尾可能有的/，然后按照/切割成段
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, p := range paths {
		// 从子节点里边找一个匹配到了当前 path 的节点
		matchChild, found := h.findMatchChild(cur, p)
		if !found {
			return nil, false
		}
		cur = matchChild
	}
	// 这里本来应该已经找完了
	if cur.handler == nil {
		// 这里适用以下场景
		// 已经注册了 /user/friends
		// 然后访问 /user
		return nil, false
	}
	return cur.handler, true
}

// Route 就相当于往树中插入节点
func (h *HandlerBasedOnTree) Route(method string, pattern string,
	handlerFunc func(c *Context)) {
	// panic("implement me")

	// 将pattern 按照URL的分隔符分割
	// 例如 /user/friends -> [user, friends]
	// 将前后的"/"去掉，统一格式
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	// 指向根节点
	cur := h.root
	for idx, path := range paths {
		// 从子节点中找一个匹配到当前 path 的节点
		mathChild, found := h.findMatchChild(cur, path)
		if found {
			cur = mathChild
		} else {
			h.createSubTree(cur, paths[idx:], handlerFunc)
			break
		}
	}

	// 离开了循环，说明加入的是短路径
	// 如先加入 /user/friends 再加入/user
	cur.handler = handlerFunc

}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string) (*node, bool) {
	for _, child := range root.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerBasedOnTree) createSubTree(root *node, paths []string, handlerFunc func(c *Context)) {
	cur := root
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}
