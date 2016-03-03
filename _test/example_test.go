package _test

import (
	"github.com/gwuhaolin/pong"
	"net/http"
	"golang.org/x/net/http2"
	"testing"
)

func Test_Hello_World(t *testing.T) {
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
	http.ListenAndServe(":3000", po)
}

func Test_HTTP2(t *testing.T) {
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
	server.ListenAndServe()
}

func Test_HTTPS(t *testing.T) {
	po := pong.New()

	// visitor https://127.0.0.1:3000/hi will see string "hi"
	po.Root.Get("/hi", func(c *pong.Context) {
		c.Response.String("hi")
	})

	http.ListenAndServeTLS(":433", "cert.pem", "key.pem", nil)
}

func TestMulti_Server(t *testing.T) {
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
		http.ListenAndServe(":3000", po0)
	}()
	http.ListenAndServe(":3001", po1)
}