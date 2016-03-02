package pong

import "net/http"

// Context represents context for the current request. It holds request and
// response objects, path parameters, data and registered handler.
// Context is handle by middleware list in order
type Context struct {
	pong      *Pong
	dataStore map[string]interface{}
	//HTTP Session
	Session *Session
	//HTTP Request,used to get params like query post-form post-file...
	Request *Request
	//HTTP Response,used to send response to client.Can send JSON XML string file...
	Response *Response
}

func newContext(pong *Pong, writer http.ResponseWriter, request *http.Request) *Context {
	context := &Context{
		pong:      pong,
		dataStore: make(map[string]interface{}),
		Request: &Request{
			HTTPRequest:     request,
			requestParamMap: make(map[string]string),
		},
		Response: &Response{
			HTTPResponseWriter: writer,
			StatusCode:         http.StatusOK,
		}}
	context.Response.context = context
	return context
}

//get a value which is set by Context.Set() method.
//if the give name is not store a nil will return
func (c *Context) Get(name string) interface{} {
	return c.dataStore[name]
}

//set a value to this context in a handle,and in next handle you can read the value by Context.Get()
//the data is store with type map[string]interface{} in memory,so set a same can overwrite old value
func (c *Context) Set(name string, value interface{}) {
	c.dataStore[name] = value
}
