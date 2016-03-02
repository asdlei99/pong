package pong

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

const (
	// Charset
	charsetUTF8 = ";charset=utf-8"

	// Headers
	httpHeaderContentType = "Content-Type"

	// Media types
	applicationJSON                  = "application/json"
	applicationJSONCharsetUTF8       = applicationJSON + charsetUTF8
	applicationJavaScript            = "application/javascript"
	applicationJavaScriptCharsetUTF8 = applicationJavaScript + charsetUTF8
	applicationXML                   = "application/xml"
	applicationXMLCharsetUTF8        = applicationXML + charsetUTF8
	textHTML                         = "text/html"
	textHTMLCharsetUTF8              = textHTML + charsetUTF8
	textPlain                        = "text/plain"
	textPlainCharsetUTF8             = textPlain + charsetUTF8
	applicationForm                  = "application/x-www-form-urlencoded"
	multipartForm                    = "multipart/form-data"
)

// is used by an HTTP handler to response to client's request.
type Response struct {
	context *Context
	// point to http.ResponseWriter in golang's standard lib
	HTTPResponseWriter http.ResponseWriter
	// HTTP status code response to client
	StatusCode int
}

// write a HTTP Header to response
//
// use before response has send to client
func (res *Response) Header(name string, value string) {
	res.HTTPResponseWriter.Header().Set(name, value)
}

// write a HTTP Cookie to response
//
// use before response has send to client
func (res *Response) Cookie(cookie *http.Cookie) {
	http.SetCookie(res.HTTPResponseWriter, cookie)
}

func (res *Response) sendData(contentType string, bs []byte) {
	res.HTTPResponseWriter.Header().Set(httpHeaderContentType, contentType)
	for _, handle := range res.context.pong.tailMiddlewareList {
		handle(res.context)
	}
	res.HTTPResponseWriter.WriteHeader(res.StatusCode)
	res.HTTPResponseWriter.Write(bs)
}

// send JSON response to client
//
// parse data by standard lib's json.Marshal and then send to client
//
// if json.Marshal fail will call HTTPErrorHandle with error and context,to handle error you should define your pong.HTTPErrorHandle
func (res *Response) JSON(data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	} else {
		res.sendData(applicationJSONCharsetUTF8, bs)
	}
}

// send JSONP response to client
//
// parse data by standard lib's json.Marshal and then send to client
// will wrap json to JavaScript's call method with give callback param
//
// if json.Marshal fail will call HTTPErrorHandle with error and context,to handle error you should define your pong.HTTPErrorHandle
func (res *Response) JSONP(data interface{}, callback string) {
	bs, err := json.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	} else {
		bs = append(append([]byte(callback), '('), append(bs, ')')...)
		res.sendData(applicationJavaScriptCharsetUTF8, bs)
	}
}

// send XML response to client
//
// parse data by standard lib's xml.Marshal and then send to client
//
// if xml.Marshal fail will call HTTPErrorHandle with error and context,to handle error you should define your pong.HTTPErrorHandle
func (res *Response) XML(data interface{}) {
	bs, err := xml.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	} else {
		res.sendData(applicationXMLCharsetUTF8, bs)
	}
}

// send a file response to client
//
// replies to the request with the contents of the named
// file or directory.
// If the provided file or direcory name is a relative path, it is
// interpreted relative to the current directory and may ascend to parent
// directories. If the provided name is constructed from user input, it
// should be sanitized before calling ServeFile. As a precaution, ServeFile
// will reject requests where r.URL.Path contains a ".." path element.
//
// As a special case, ServeFile redirects any request where r.URL.Path
// ends in "/index.html" to the same path, without the final
// "index.html". To avoid such redirects either modify the path or
// use ServeContent.
func (res *Response) File(filePath string) {
	http.ServeFile(res.HTTPResponseWriter, res.context.Request.HTTPRequest, filePath)
}

// send String response to client
func (res *Response) String(str string) {
	res.sendData(textPlainCharsetUTF8, []byte(str))
}

// send HTML response to client by html string
func (res *Response) HTML(html string) {
	res.sendData(textHTMLCharsetUTF8, []byte(html))
}

// send HTML response to client by render HTML template with give data
//
// LoadTemplateGlob before use Render
func (res *Response) Render(template string, data interface{}) {
	tpl := res.context.pong.htmlTemplate
	if tpl != nil {
		html := bytes.Buffer{}
		err := tpl.ExecuteTemplate(&html, template, data)
		if err != nil {
			res.context.pong.HTTPErrorHandle(err, res.context)
		} else {
			res.sendData(textHTMLCharsetUTF8, html.Bytes())
		}
	} else {
		fmt.Errorf("pong:LoadTemplateGlob before use Render")
		res.context.pong.NotFindHandle(res.context)
	}
}

// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
//
// The Response.StatusCode should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
func (res *Response) Redirect(url string) {
	http.Redirect(res.HTTPResponseWriter, res.context.Request.HTTPRequest, url, res.StatusCode)
}
