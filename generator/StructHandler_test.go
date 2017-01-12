package generator

import (
	"testing"
)

func TestStructHandler(t *testing.T) {
	var structure = StructHandler{
		Dir:    "./test/handler",
		Name:   "Index",
		Method: "POST",
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
	t.Log(structure.Output())
}
