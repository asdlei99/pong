# [pong](https://github.com/gwuhaolin/pong)

[![Build Status](https://travis-ci.org/gwuhaolin/pong.svg?branch=master)](https://travis-ci.org/gwuhaolin/pong)
[![Coverage Status](https://coveralls.io/repos/github/gwuhaolin/pong/badge.svg?branch=master)](https://coveralls.io/github/gwuhaolin/pong?branch=master)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gwuhaolin/pong#SessionManager)

A router for high performance web service write in golang.

# Introduction
Pong is just a http router library, use it make high performance web service in minutes.
It's work is to route a request to register handle then provide convenient function to get param from request and send response and also provide option **HTTP session** support.
Pong process every request as a product in production line, use register **middleware** do some change to the product.This like the way in NodeJs's famous Express do.
It has **no dependency** small and clear, support **route conflict tips**.

# Performance

# Hello World
```go
    package main
    import (
    	"github.com/gwuhaolin/pong"
    	"net/http"
    )
    func main() {
    	po := pong.New()
        root := po.Root
    	// visit http://127.0.0.1:3000/ping will see string "pong"
    	root.Get("/ping", func(c *pong.Context) {
    		c.Response.String("pong")
    	})

    	// a sub router
    	sub := root.Router("/sub")

    	// visit http://127.0.0.1:3000/sub/pong will see JSON "{"name":"pong"}"
    	sub.Get("/:name", func(c *pong.Context) {
    		m := map[string]string{
    			"name":c.Request.Param("name"),
    		}
    		c.Response.JSON(m)
    	})

    	// Run Server Listen on 127.0.0.1:3000
    	http.ListenAndServe(":3000", po)
    }
```

# Installation
```bash
    go get github.com/gwuhaolin/pong
```

# Principle

# Catalogue

# Listen and Server
pong not provide Listen and Server, it just do thing about route and handle, so you can should standard lib's function
### HTTPS
```go
    po := pong.New()

	// visit https://127.0.0.1:3000/hi will see string "hi"
	root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	http.ListenAndServeTLS(":433", "cert.pem", "key.pem", nil)
```
### HTTP2
```go
    po := pong.New()

	// visit http://127.0.0.1:3000/hi will see string "hi"
	root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	server := &http.Server{
		Handler:po,
		Addr:":3000",
	}
	http2.ConfigureServer(server, &http2.Server{})
	server.ListenAndServe()
```
### Multi-Server
```go
    po0 := pong.New()
	po1 := pong.New()

	// visit http://127.0.0.1:3000/hi will see string "0"
	po0.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("0")
	})

	// visit http://127.0.0.1:3001/hi will see string "1"
	po1.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("1")
	})
	go func() {
		http.ListenAndServe(":3000", po0)
	}()
	http.ListenAndServe(":3001", po1)
```

# Route
Route every request to the right handle is pong's job.
Pong will build a tree in type map when you register your handle to a path, when server has run and request come in, pong will use map's hash algorithm to find a register handle.
Pong's router not support regular expression because infrequency and avoid it can improve performance
Pong support sub Router, a route like a tree which is comprise by one or more sub Router
Pong's Root Router can access by `pong.Root` which point to root path `/`
### HTTP Methods
after route a request to path, pong can also route diff HTTP method. This `Delete` `Get` `Head` `Options` `Patch` `Post` `Put` `Trace` `Any` are support.
```go
    root := po.Root
    root.Delete("/", func(c *Context) {
		c.Response.String("Delete")
	})
    root.Put("/", func(c *Context) {
		c.Response.String("Put")
	})
    root.Any("/", func(c *Context) {
		c.Response.String("Any will overwrite all of them because is registed last, this means overwrite by registed order")
	})
```
### Sub Router
```go
	// visit / will see string "/"
    root.Get("/", func(c *Context) {
		c.Response.String("/")
	})
	sub := root.Router("sub")
	// visit /sub/hi will see string "sub"
	sub.Get("/hi", func(c *Context) {
        c.Response.String("sub")
    })
    sub2 := sub.Router("sub")
	// visit /sub/sub2/hi will see string "sub2"
	sub2.Get("/hi", func(c *Context) {
        c.Response.String("sub2")
    })
```
### Path Param
```go
	// visit /abc will see string "abc"
    root.Get("/:param", func(c *Context) {
		c.Response.String(c.Request.Param("param"))
	})
	// param in router path
	sub := root.Router("sub/:name")
	// visit /sub/abc/hi will see string "abc"
	sub.Get("/hi", func(c *Context) {
        c.Response.String(c.Request.Param("param"))
    })
```
### WebSocket
TODO
### Route Conflict Tips
see Route Conflict this code:
```go
	root.Get("/path", func(c *Context) {
		c.Response.String("path")
	})
	root.Get("/:name", func(c *Context) {
		c.Response.String(c.Request.Param("name"))
	})
```
`:name` march `/path`, when this happen pong will print warning to tell developer fix Conflict. But this code can still run, pong has rule you must know:
**path's(/path) priority level is high than param's(/:name)**, so for this code when you:
- visit `/path` will see string `path`, which use handle set in `root.Get("/path",handle)`
- visit `/hal` will see string `hal`, which use handle set in `root.Get("/:name",handle)`

# Request
### Query Param
```go
	// visit /?param=abc will see string "abc"
    root.Get("/", func(c *Context) {
		c.Response.String(c.Request.Query("param"))
	})
```
### Form Param
```go
	// post / with body "param=abc" will see string "abc"
    root.Post("", func(c *Context) {
		c.Response.String(c.Request.Form("param"))
	})
```
### Post File
```go
    // post / with a text file will see file's context
    root.Post("/", func(c *Context) {
		file, _, _ := c.Request.File("file")
		bs, _ := ioutil.ReadAll(file)
		c.Response.String(string(bs))
	})
```
### Bind
Pong provide convenient way to parse request's params and bind to a struct
### BindJSON
parse request's body data as JSON and use standard lib json.Unmarshal to bind data to struct
```go
    type testUser struct {
    	Name  string
    	Age   int
    }
    // post / with a json string will see json again
    root.Post("/", func(c *Context) {
		user := testUser{}
		c.Request.BindJSON(&bindUser)
		c.Response.JSON(user)
	})
```
### BindXML
parse request's body data as XML and use standard lib XML.Unmarshal to bind data to struct
```go
    // post / with a xml string will see xml again
    root.Post("/", func(c *Context) {
		user := testUser{}
		c.Request.BindXML(&bindUser)
		c.Response.XML(user)
	})
```
### BindForm
parse request's body post form as map and bind data to struct use filed name
```go
    // post / with a name=hal&age=23 will see json "{"name":"hal","age":23}"
    root.Post("/", func(c *Context) {
		user := testUser{}
		c.Request.BindXML(&bindUser)
		c.Response.JSON(user)
	})
```
### BindQuery
parse request's query params as map and bind data to struct use filed name
```go
    // visit /?name=hal&age=23 will see json "{"name":"hal","age":23}"
    root.Post("/", func(c *Context) {
		user := testUser{}
		c.Request.BindQuery(&bindUser)
		c.Response.JSON(user)
	})
```
### AutoBind
auto bind will look request's http Header `ContentType`
- if request ContentType is applicationJSON will use `BindJSON` to parse
- if request ContentType is applicationXML will use `BindXML` to parse
- if request ContentType is applicationForm or multipartForm will use `BindForm` to parse
- else will return an ErrorTypeNotSupport error
```go
    // post / with a json "{"name":"hal","age":23}" will see "{"name":"hal","age":23}"
    // post / with a name=hal&age=23 will see json "{"name":"hal","age":23}"
    // visit /?name=hal&age=23 will see json "{"name":"hal","age":23}"
    // post / with a xml will see "{"name":"hal","age":23}"
    root.Post("/", func(c *Context) {
		user := testUser{}
		c.Request.AutoBind(&bindUser)
		c.Response.JSON(user)
	})
```
# Response
### Set Header
write a HTTP Header to response use before response has send to client
```go
    c.Response.Header("X-name", "mine header")
```
### Set Cookie
write a HTTP Cookie to response
```go
    c.Response.Cookie(&http.Cookie{Name: "id", Value: "123"})
```
### Send JSON
send JSON response to client, parse data by standard lib's json.Marshal and then send to client
### Send JSONP
parse data by standard lib's json.Marshal and then send to client will wrap json to JavaScript's call method with give callback param
```go
    // visit /hi will see json "callback({"name":"hal","age":23})"
    root.Get("/hi", func(c *Context) {
        user := testUser{
                Name:"hal",
                Age:23,
        }
		c.Response.JSONP(user,"callback")
	})
```
### Send XML
parse data by standard lib's xml.Marshal and then send to client
### Send File
send a file response to client
```go
    // visit /hi will see file hi.zip
    root.Get("/hi", func(c *Context) {
		c.Response.File("hi.zip")
	})
```
### Send String
send String response to client
### Redirect
Redirect replies to the request with a redirect to url, which may be a path relative to the request path.
```go
    // visit /redirect will redirect to /
    root.Get("/redirect", func(c *Context) {
		c.Response.Redirect("/")
	})
```
### Render HTML Template
send HTML response to client by render HTML template with give data, LoadTemplateGlob before use Render
```go
    po.LoadTemplateGlob("*.html")
    // visit /index will see index.html template render by data
    root.Get("/:name", func(c *Context) {
    		c.Response.Render(name, dataToRender)
    })
````

# Middleware
### Router Middleware
### Tail Middleware
### Log HTTP Request
### Log HTTP Request To MongoDB
### Write Your Middleware

# Session
### Set and Get
### Reset
### Destory
### Store Session In Redis
### Store Session In MongoDB
### Store Session In SQL
### Write Your Session Manager

# LICENSE
The MIT License (MIT)
Copyright (c) 2016 吴浩麟
