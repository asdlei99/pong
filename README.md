# [pong](https://github.com/gwuhaolin/pong)

[![Build Status](https://travis-ci.org/gwuhaolin/pong.svg?branch=master)](https://travis-ci.org/gwuhaolin/pong)
[![Coverage Status](https://coveralls.io/repos/github/gwuhaolin/pong/badge.svg?branch=master)](https://coveralls.io/github/gwuhaolin/pong?branch=master)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gwuhaolin/pong#SessionManager)

A simple HTTP router for golang.

# Introduction
Pong is just a http router library.
It's work is to router a request to register handle then provide convenient function to get param from request and send response and also provide option HTTP session support.
Pong process every request as a product in production line, use register middleware do some change to the product.This like the way in NodeJs's famous Express do.
It's api is small and clear, no dependency, good performance.

# Performance

# Hello World
```
    func main() {
    	po := pong.New()

    	// visitor http://127.0.0.1:3000/hi will see string "hi"
    	po.Root.Get("/hi", func(c *Context) {
		    c.Response.String("/hi")
	    })

	    // a sub router
	    sub := po.Root.Router("/sub")

	    // visitor http://127.0.0.1:3000/sub/pong will see string "hello pong"
	    sub.Get("/:name", func(c *Context) {
		    c.Response.String("hello " + c.Request.Param("name"))
	    })

	    // Run Server Listen on 127.0.0.1:3000
        http.ListenAndServe(":3000", po)
   }
```

# Installation
```
go get github.com/gwuhaolin/pong
```

# Principle

# Catalogue

# Server
## Listen Address
## HTTPS
## HTTP2
## Multi-Server

# Route
## Path Param
## Sub Router

# Request
## Path Param
## Query Param
## Form Param
## Post File
## Bind
## BindJSON
## BindXML
## BindForm
## BindForm
## BindQuery
## AutoBind

# Response
## Set Header
## Set Cookie
## Send JSON
## Send JSONP
## Send XML
## Send File
## Send String
## Redirect
## Render HTML Template

# Middleware
## Router Middleware
## Tail Middleware
## Log HTTP Request
## Log HTTP Request To MongoDB
## Write Your's Middleware

# Session
## Set and Get
## Reset
## Destory
## Store Session In Redis
## Store Session In MongoDB
## Store Session In SQL
## Write Your's Session Manager