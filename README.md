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
```go
    package main

    import (
    	"github.com/gwuhaolin/pong"
    	"net/http"
    	"log"
    )

    func main() {
    	po := pong.New()

    	// visitor http://127.0.0.1:3000/ping will see string "pong"
    	po.Root.Get("/ping", func(c *pong.Context) {
    		c.Response.String("pong")
    	})

    	// a sub router
    	sub := po.Root.Router("/sub")

    	// visitor http://127.0.0.1:3000/sub/pong will see JSON "{"name":"pong"}"
    	sub.Get("/:name", func(c *pong.Context) {
    		m := map[string]string{
    			"name":c.Request.Param("name"),
    		}
    		c.Response.JSON(m)
    	})

    	// Run Server Listen on 127.0.0.1:3000
    	log.Println(http.ListenAndServe(":3000", po))
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

	// visitor https://127.0.0.1:3000/hi will see string "hi"
	po.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	log.Fatal(http.ListenAndServeTLS(":433", "cert.pem", "key.pem", nil))
```
### HTTP2
```go
    po := pong.New()

	// visitor http://127.0.0.1:3000/hi will see string "hi"
	po.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	server := &http.Server{
		Handler:po,
		Addr:":3000",
	}
	http2.ConfigureServer(server, &http2.Server{})
	log.Fatal(server.ListenAndServe())
```
### Multi-Server
```go
    po0 := pong.New()
	po1 := pong.New()

	// visitor https://127.0.0.1:3000/hi will see string "0"
	po0.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("0")
	})

	// visitor https://127.0.0.1:3001/hi will see string "1"
	po1.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("1")
	})
	go func() {
		log.Fatal(http.ListenAndServe(":3000", po0))
	}()
	log.Fatal(http.ListenAndServe(":3001", po1))
```

# Route
### Path Param
### Sub Router
### WebSocket

# Request
### Path Param
### Query Param
### Form Param
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
### Write Your's Middleware

# Session
### Set and Get
### Reset
### Destory
### Store Session In Redis
### Store Session In MongoDB
### Store Session In SQL
### Write Your's Session Manager

# LICENSE
The MIT License (MIT) Copyright (c) 2016 吴浩麟
