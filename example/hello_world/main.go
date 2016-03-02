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
	log.Fatal(http.ListenAndServe(":3000", po))
}