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

    	// visit http://127.0.0.1:3000/ping will see string "pong"
    	po.Root.Get("/ping", func(c *pong.Context) {
    		c.Response.String("pong")
    	})

    	// a sub router
    	sub := po.Root.Router("/sub")

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
	po.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	http.ListenAndServeTLS(":433", "cert.pem", "key.pem", nil)
```
### HTTP2
```go
    po := pong.New()

	// visit http://127.0.0.1:3000/hi will see string "hi"
	po.Root.Get("/hi", func(c *pong.Context) {
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
Pong will build a tree in type map when you register your handle to path, when server has run and request come in, pong will use map's hash algorithm to find a register handle.
Pong's router not support regular expression because infrequency and avoid it can improve performance
Pong support sub Router, a route like a tree which is comprise by one or more sub Router
Pong's Root Router can access by pong.Root which point to root path /
### HTTP Methods
after route a request to path, pong can also route diff HTTP method. This `Delete` `Get` `Head` `Options` `Patch` `Post` `Put` `Trace` `Any` are support.
```go
    po.Root.Delete("/", func(c *Context) {
		c.Response.String("Delete")
	})
    po.Root.Put("/", func(c *Context) {
		c.Response.String("Put")
	})
    po.Root.Any("/", func(c *Context) {
		c.Response.String("Any will overwrite all of them because is registed last, this means overwrite by registed order")
	})
```
### Sub Router
```go
	// visit / will see string "/"
    po.Root.Get("/", func(c *Context) {
		c.Response.String("/")
	})
	sub := po.Root.Router("sub")
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
    po.Root.Get("/:param", func(c *Context) {
		c.Response.String(c.Request.Param("param"))
	})
	// param in router path
	sub := po.Root.Router("sub/:name")
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
	po.Root.Get("/path", func(c *Context) {
		c.Response.String("path")
	})
	po.Root.Get("/:name", func(c *Context) {
		c.Response.String(c.Request.Param("name"))
	})
```
`:name` march `/path`, when this happen pong will print warning to tell developer fix Conflict. But this code can still run, pong has rule you must know:
**path's(/path) priority level is high than param's(/:name)**, so for this code when you visit:
- `/path` will see string `path`, which use handle set in `po.Root.Get("/path",handle)`
- `/hal` will see string `hal`, which use handle set in `po.Root.Get("/:name",handle)`

# Request
### Query Param
```go
	// visit /?param=abc will see string "abc"
    po.Root.Get("/", func(c *Context) {
		c.Response.String(c.Request.Query("param"))
	})
```
### Form Param
```go
	// post / with body "param=abc" will see string "abc"
    po.Root.Post("", func(c *Context) {
		c.Response.String(c.Request.Form("param"))
	})
```
### Post File
### Bind
### BindJSON
### BindXML
### BindForm
### BindForm
### BindQuery
### AutoBind

# Response
### Set Header
### Set Cookie
### Send JSON
### Send JSONP
### Send XML
### Send File
### Send String
### Redirect
### Render HTML Template

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
