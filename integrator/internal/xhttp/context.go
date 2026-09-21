package xhttp

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Token    string
	UserKey  string
	IsBearer bool
}

func New(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) Path() string {
	return c.Request.URL.Path
}

func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) Header(name string) string {
	return c.Request.Header.Get(name)
}

func (c *Context) SetHeader(name string, content string) {
	c.Writer.Header().Set(name, content)
}

func (c *Context) Post(name string) string {
	_ = c.Request.ParseForm()
	return c.Request.Form.Get(name)
}

func (c *Context) JSON(code int, data any) {
	c.SetHeader("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	_ = json.NewEncoder(c.Writer).Encode(data)
}

func (c *Context) Success(data any) {
	c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"result":  data,
	})
}

func (c *Context) Error(code int, message string) {
	c.JSON(code, map[string]any{
		"success": false,
		"message": message,
	})
}

func (c *Context) JSONBody(dest any) error {
	return json.NewDecoder(c.Request.Body).Decode(dest)
}
