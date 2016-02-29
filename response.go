package pong

import (
	"net/http"
	"encoding/json"
	"encoding/xml"
	"bytes"
)

const (
// Charset
	charsetUTF8 = ";charset=utf-8"

// Headers
	httpHeaderContentType = "Content-Type"

// Media types
	applicationJSON = "application/json"
	applicationJSONCharsetUTF8 = applicationJSON + charsetUTF8
	applicationJavaScript = "application/javascript"
	applicationJavaScriptCharsetUTF8 = applicationJavaScript + charsetUTF8
	applicationXML = "application/xml"
	applicationXMLCharsetUTF8 = applicationXML + charsetUTF8
	textHTML = "text/html"
	textHTMLCharsetUTF8 = textHTML + charsetUTF8
	textPlain = "text/plain"
	textPlainCharsetUTF8 = textPlain + charsetUTF8
	applicationForm = "application/x-www-form-urlencoded"
	multipartForm = "multipart/form-data"
)

type Response struct {
	context            *Context
	HTTPResponseWriter http.ResponseWriter
	StatusCode         int
}

func (res *Response)Header(name string, value string) {
	res.HTTPResponseWriter.Header().Set(name, value)
}

func (res *Response)Cookie(cookie *http.Cookie) {
	http.SetCookie(res.HTTPResponseWriter, cookie)
}

func (res *Response)sendData(contentType string, bs []byte) {
	res.HTTPResponseWriter.Header().Set(httpHeaderContentType, contentType)
	for _, handle := range res.context.pong.tailMiddlewareList {
		handle(res.context)
	}
	res.HTTPResponseWriter.WriteHeader(res.StatusCode)
	_, err := res.HTTPResponseWriter.Write(bs)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	}
}

func (res *Response)JSON(data interface{}) {
	bs, err := json.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	}else {
		res.sendData(applicationJSONCharsetUTF8, bs)
	}
}

func (res *Response)JSONP(data interface{}, callback string) {
	bs, err := json.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	}else {
		bs = append(append([]byte(callback), '('), append(bs, ')')...)
		res.sendData(applicationJavaScriptCharsetUTF8, bs)
	}
}

func (res *Response)XML(data interface{}) {
	bs, err := xml.Marshal(data)
	if err != nil {
		res.context.pong.HTTPErrorHandle(err, res.context)
	}else {
		res.sendData(applicationXMLCharsetUTF8, bs)
	}
}

func (res *Response)File(filePath string) {
	http.ServeFile(res.HTTPResponseWriter, res.context.Request.HTTPRequest, filePath)
}

func (res *Response)String(str string) {
	res.sendData(textPlainCharsetUTF8, []byte(str))
}

func (res *Response)HTML(html string) {
	res.sendData(textHTMLCharsetUTF8, []byte(html))
}

func (res *Response)Render(template string, data interface{}) {
	tpl := res.context.pong.htmlTemplate
	if tpl != nil {
		html := bytes.Buffer{}
		err := tpl.ExecuteTemplate(&html, template, data)
		if err != nil {
			res.context.pong.HTTPErrorHandle(err, res.context)
		}else {
			res.sendData(textHTMLCharsetUTF8, html.Bytes())
		}
	}else {
		panic("LoadTemplateGlob before use Render")
	}
}

func (res *Response)Redirect(url string) {
	http.Redirect(res.HTTPResponseWriter, res.context.Request.HTTPRequest, url, res.StatusCode)
}
