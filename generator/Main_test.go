package generator

import (
	"testing"
)

func TestMain1(t *testing.T) {
	var structure = &StructHandler{
		Dir:     "./test/handler",
		Name:    "Index",
		UrlPath: "/test/index",
		Method:  "POST",
		Fields: []Field{
			{
				Type: "string",
				Name: "Title",
				In:   "query",
				Desc: "title param",
			},
			{
				Type: "*multipart.FileHeader",
				Name: "Img",
				In:   "formData",
				Desc: "img param",
			},
			{
				Type: "*http.Cookie",
				Name: "Cookie",
				In:   "cookie",
				Desc: "cookie param",
			},
		},
		Note:   "index handler",
		Return: "{}",
	}
	var function = &FuncHandler{
		Dir:     "./test/handler",
		Name:    "Index2",
		UrlPath: "/test/index2",
		Method:  "GET",
	}
	var middleware = &FuncHandler{
		Dir:     "./test/middleware",
		Name:    "middleware",
		UrlPath: "/test",
	}

	var router, err = NewRouter("MyappRoute", "./test/router")
	if err != nil {
		t.Logf("%v", err)
	}
	router.AddHandler(structure)
	router.AddHandler(function)
	router.AddMiddleware(middleware)
	router.AddStatic("static fs", "/test/static", "./test/static/test", true, false)
	router.AddStatic("static fs", "/test2/static", "./test/static/test", true, false)

	m, err := NewMain("./test")
	if err != nil {
		t.Logf("%v", err)
	}
	m.AddFrame(router, "myapp", "1.0")
	t.Log(m.Output())
}

func TestMain2(t *testing.T) {
	var structure = &StructHandler{
		Dir:     "./test",
		Name:    "Index",
		UrlPath: "/test/index",
		Method:  "POST",
		Fields: []Field{
			{
				Type: "string",
				Name: "Title",
				In:   "query",
				Desc: "title param",
			},
			{
				Type: "*multipart.FileHeader",
				Name: "Img",
				In:   "formData",
				Desc: "img param",
			},
			{
				Type: "*http.Cookie",
				Name: "Cookie",
				In:   "cookie",
				Desc: "cookie param",
			},
		},
		Note:   "index handler",
		Return: "{}",
	}

	var router, err = NewRouter("MyappRoute", "./test")
	if err != nil {
		t.Logf("%v", err)
	}
	router.AddHandler(structure)
	router.AddStatic("static fs", "/test/static", "./test/static/test", true, false)

	m, err := NewMain("./test")
	if err != nil {
		t.Logf("%v", err)
	}
	m.AddFrame(router, "myapp")
	t.Log(m.Output())
}
