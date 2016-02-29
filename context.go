package pong

import "net/http"

type Context struct {
	pong      *Pong
	dataStore map[string]interface{}
	Session   *Session
	Request   *Request
	Response  *Response
}

func newContext(pong *Pong, writer http.ResponseWriter, request  *http.Request) *Context {
	context := &Context{
		pong:pong,
		dataStore:make(map[string]interface{}),
		Request:&Request{
			HTTPRequest:request,
			requestParamMap:make(map[string]string),
		},
		Response:&Response{
			HTTPResponseWriter:writer,
			StatusCode:http.StatusOK,
		}}
	context.Response.context = context
	return context
}

func (c *Context)Get(name string) interface{} {
	return c.dataStore[name]
}

func (c *Context)Set(name string, value interface{}) {
	c.dataStore[name] = value
}