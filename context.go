package cill

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// base info
	Writer http.ResponseWriter
	Req    *http.Request
	engine *Engine
	// request info
	Path   string
	Method string
	Params map[string]string // dynamic router params
	// response info
	StatusCode int
	// handler
	handler []HandlerFunc
	idx int
}

func NewContext(w http.ResponseWriter, req *http.Request, e *Engine) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		engine: e,
		idx: -1,
	}
}

/*
	Context handler function
*/
func (c *Context) Next() {
	c.idx ++
	n := len(c.handler)
	for ;c.idx < n; c.idx ++ {
		c.handler[c.idx](c)
	}
}


/*
	Context request function
*/

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

/*
	Context response function
*/

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, format string, values ...any) error {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	return nil
}

func (c *Context) JSON(code int, obj any) error {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encode := json.NewEncoder(c.Writer)
	if err := encode.Encode(obj); err != nil {
		return fmt.Errorf("encode json fail. %v", err)
	}
	return nil
}

func (c *Context) Data(code int, data []byte) error {
	c.Status(code)
	if err := c.write(data); err != nil {
		return err
	}
	return nil
}

func (c *Context) HTML(code int, html string) error {
	c.Status(code)
	if err := c.write([]byte(html)); err != nil {
		return err
	}
	return nil
}

func (c *Context) write(data []byte) error {
	n, err := c.Writer.Write(data)
	if n != len(data) {
		return fmt.Errorf("write data to response fail. write length is %d, want length is %d", n, len(data))
	}
	if err != nil {
		return fmt.Errorf("write data to response fail. %v", err)
	}
	return nil
}
