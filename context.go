package pong

import "net/http"

// Incoming requests to a server will create a Context
// Context can be use to:
// 	get HTTP request param query...
// 	send HTTP response like JSON JSONP XML HTML file ...
// Context is handle by middleware list in order
type Context struct {
	pong      *Pong
	dataStore map[string]interface{}
	//
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