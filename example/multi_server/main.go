package main

import (
	"github.com/gwuhaolin/pong"
	"net/http"
	"log"
)

func main() {
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
}