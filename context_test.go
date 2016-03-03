package pong

import (
	"net/http"
	"reflect"
	"testing"
	"github.com/gwuhaolin/pong/_test"
)

func TestContext(t *testing.T) {
	po, baseURL := runPong()
	root := po.Root
	user := _test_util.TestUser{
		Name:  "吴浩麟",
		Age:   23,
		Money: 123.456,
		Alive: true,
		Notes: []_test_util.TestNote{
			{Text: "明天去放风筝"},
			{Text: "今天我们去逛宜家啦"},
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
		t.Log(`TestContext`)
	})
	defer http.Get(baseURL + "/user")
}
