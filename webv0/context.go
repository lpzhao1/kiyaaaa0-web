package webv0

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(bs)
	if err != nil {
		return err
	}
	c.W.WriteHeader(status)
	return nil
}
