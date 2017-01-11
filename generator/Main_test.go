package generator

import (
	"testing"
)

func TestMain1(t *testing.T) {
	var structure = &StructHandler{
		Dir:     "./test/handler",
		Name:    "Index",
		UrlPath: "/test/index",
		Method:  "GET",
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

	var router, err = NewRouter("route", "./test/router")
	if err != nil {
		t.Logf("%v", err)
	}
	router.API(structure)

	m, err := NewMain("./test")
	if err != nil {
		t.Logf("%v", err)
	}
	m.AddFrame(router, "myapp", "1.0")
	t.Log(m.CreateFile())
}

func TestMain2(t *testing.T) {
	var structure = &StructHandler{
		Dir:     "./test",
		Name:    "Index",
		UrlPath: "/test/index",
		Method:  "GET",
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

	var router, err = NewRouter("route", "./test")
	if err != nil {
		t.Logf("%v", err)
	}
	router.API(structure)

	m, err := NewMain("./test")
	if err != nil {
		t.Logf("%v", err)
	}
	m.AddFrame(router, "myapp")
	t.Log(m.CreateFile())
}
