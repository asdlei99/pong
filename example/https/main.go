package main

import (
	"github.com/gwuhaolin/pong"
	"net/http"
	"log"
)

func main() {
	po := pong.New()

	// visitor https://127.0.0.1:3000/hi will see string "hi"
	po.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	log.Fatal(http.ListenAndServeTLS(":433", "cert.pem", "key.pem", nil))
}