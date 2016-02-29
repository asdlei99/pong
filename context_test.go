package pong

import (
	"net/http"
	"testing"
	"reflect"
)

func TestContext(t *testing.T) {
	po := New()
	go http.ListenAndServe(listenAddr, po)
	po.Root.Middleware(logRequest)
	root := po.Root
	user := testUser{
		Name:"吴浩麟",
		Age:23,
		Money:123.456,
		Alive:true,
		Notes:[]testNote{
			{Text:"明天去放风筝"},
			{Text:"今天我们去逛宜家啦"},
		},
	}
	root.Middleware(func(c *Context) {
		c.Set("name", user)
	})
	root.Get("/user", func(c *Context) {
		resUser := c.Get("name")
		if !reflect.DeepEqual(resUser, user) {
			t.Error(resUser)
		}
	})
	defer http.Get(baseURL + "/user")
}
