package main

import (
	"github.com/gwuhaolin/pong"
	"net/http"
	"log"
	"golang.org/x/net/http2"
)

func main() {
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
}