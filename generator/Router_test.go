package generator

import (
	"testing"
)

func TestRouter(t *testing.T) {
	var structure = &StructHandler{
		Dir:     "./test/handler",
		Name:    "Index",
		Method:  "GET",
		UrlPath: "/test/index",
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
	structure.CreateFile()

	var router, err = NewRouter("route", "./test/router")
	if err != nil {
		t.Logf("%v", err)
	}
	router.API(structure)
	t.Log(router.CreateFile())
}
