// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generator

import (
	"fmt"
	"strings"

	"github.com/henrylee2cn/faygo"
)

/*Param tag value description:
    tag   |   key    | required |     value     |   desc
    ------|----------|----------|---------------|----------------------------------
    param |    in    | only one |     path      | (position of param) if `required` is unsetted, auto set it. e.g. url: "http://www.abc.com/a/{path}"
    param |    in    | only one |     query     | (position of param) e.g. url: "http://www.abc.com/a?b={query}"
    param |    in    | only one |     formData  | (position of param) e.g. "request body: a=123&b={formData}"
    param |    in    | only one |     body      | (position of param) request body can be any content
    param |    in    | only one |     header    | (position of param) request header info
    param |    in    | only one |     cookie    | (position of param) request cookie info, support: `*http.Cookie`,`http.Cookie`,`string`,`[]byte`
    param |   name   |    no    |   (e.g.`id`)   | specify request param`s name
    param | required |    no    |               | request param is required
    param |   desc   |    no    |   (e.g.`id`)   | request param description
    param |   len    |    no    | (e.g.`3:6` `3`) | length range of param's value
    param |   range  |    no    |  (e.g.`0:10`)  | numerical range of param's value
    param |  nonzero |    no    |               | param`s value can not be zero
    param |   maxmb  |    no    |   (e.g.`32`)   | when request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
    param |  regexp  |    no    | (e.g.`^\\w+$`) | verify the value of the param with a regular expression(param value can not be null)
    param |   err    |    no    |(e.g.`incorrect password format`)| the custom error for binding or validating

    NOTES:
        1. the binding object must be a struct pointer
        2. in addition to `*multipart.FileHeader`, the binding struct's field can not be a pointer
        3. `regexp` or `param` tag is only usable when `param:"type(xxx)"` is exist
        4. if the `param` tag is not exist, anonymous field will be parsed
        5. when the param's position(`in`) is `formData` and the field's type is `*multipart.FileHeader`, `multipart.FileHeader`, `[]*multipart.FileHeader` or `[]multipart.FileHeader`, the param receives file uploaded
        6. if param's position(`in`) is `cookie`, field's type must be `*http.Cookie` or `http.Cookie`
        7. param tags `in(formData)` and `in(body)` can not exist at the same time
        8. there should not be more than one `in(body)` param tag

List of supported param value types:
    base    |   slice    | special
    --------|------------|-------------------------------------------------------
    string  |  []string  | [][]byte
    byte    |  []byte    | [][]uint8
    uint8   |  []uint8   | *multipart.FileHeader (only for `formData` param)
    bool    |  []bool    | []*multipart.FileHeader (only for `formData` param)
    int     |  []int     | *http.Cookie (only for `net/http`'s `cookie` param)
    int8    |  []int8    | http.Cookie (only for `net/http`'s `cookie` param)
    int16   |  []int16   | struct (struct type only for `body` param or as an anonymous field to extend params)
    int32   |  []int32   |
    int64   |  []int64   |
    uint8   |  []uint8   |
    uint16  |  []uint16  |
    uint32  |  []uint32  |
    uint64  |  []uint64  |
    float32 |  []float32 |
    float64 |  []float64 |
*/
type (
	// StructHandler struct handler
	StructHandler struct {
		Dir          string  // file path or package name
		UrlPath      string  // URL's path
		Name         string  // (required) struct name
		Fields       []Field // fields
		ServeContent string  // main logic
		Note         string  // note for API
		Return       string  // response content demo
		Method       faygo.Methodset
		fileParams   []string
		filesParams  []string
		importmap    map[string]bool
		sign         string
		isMainPkg    bool
	}
	// Field struct handler's field
	Field struct {
		Type      string // (required) Value type
		Name      string // (required) Field name
		ParamName string // Param name for API
		// Default   string // Default value for API doc
		In        string // The position of the parameter; if empty, indicates that it is not a request parameter
		Required  bool   // Is a required parameter
		Nonzero   bool   // Param`s value can not be zero
		Len       string // Length range of param's value
		Range     string // Numerical range of param's value
		Regexp    string // Verify the value of the param with a regular expression(param value can not be null)
		Maxmb     int    // When request Content-Type is multipart/form-data, the max memory for body.(multi-param, whichever is greater)
		Err       string // The custom error for binding or validating
		Desc      string // Description
		OtherTags string // Other tags string
		isParam   bool
	}
)

// Output creates struct handler file.
func (s *StructHandler) Output() error {
	code, err := s.Create()
	if err != nil {
		return err
	}
	return writeFile(s.Dir, faygo.SnakeString(s.Name)+".go", code)
}

// Create returns struct handler's codes
func (s *StructHandler) Create() (code string, err error) {
	// initialize
	err = s.init()
	if err != nil {
		return
	}
	// build codes
	code = fmt.Sprintf("package %s\n%s\n%s", s.PkgName(), importCode(s.importmap), s.createStruct())
	return code, nil
}

// GetUrlPath returns router node's url path.
func (s *StructHandler) GetUrlPath() string {
	return s.UrlPath
}

// GetMethod returns request method.
func (s *StructHandler) GetMethod() faygo.Methodset {
	return s.Method
}

// GetName returns handler type name.
func (s *StructHandler) GetName() string {
	return s.Name
}

// RouterName returns router node's name
func (s *StructHandler) RouterName() string {
	if len(s.Note) > 0 {
		return strings.Split(s.Note, "\n")[0]
	}
	return s.Name
}

// TryMainPkg tries to set it as the main package
func (s *StructHandler) TryMainPkg(mainPkgPath string) {
	if s.Dir != mainPkgPath {
		return
	}
	s.isMainPkg = true
}

// PkgPath returns the package path, e.g `github.com/henrylee2cn/fay/test`
func (s *StructHandler) PkgPath() string {
	if s.isMainPkg || s.Dir == "" {
		return ""
	}
	dirs := strings.Split(s.Dir, "/src/")
	if len(dirs) < 2 {
		faygo.Fatalf("You must generate codes in the `src` or its offspring directory!")
	}
	return strings.Join(dirs[1:], "/src/")
}

// PkgName returns the package name, e.g `handler`
func (s *StructHandler) PkgName() string {
	if s.isMainPkg || s.Dir == "" {
		return "main"
	}
	return s.Dir[strings.LastIndex(s.Dir, "/")+1:]
}

// PkgPrefix returns the package name, e.g `handler.`
func (s *StructHandler) PkgPrefix() string {
	if s.isMainPkg || s.Dir == "" {
		return ""
	}
	return s.Dir[strings.LastIndex(s.Dir, "/")+1:] + "."
}

// initialize
func (s *StructHandler) init() error {
	// if len(s.Method.Methods()) == 0 {
	// 	return errors.New("The Method field must be setted correctly.")
	// }
	err := cleanDir(&s.Dir)
	if err != nil {
		return err
	}
	err = cleanUrlPath(&s.UrlPath)
	if err != nil {
		return err
	}
	s.sign, err = cleanName(&s.Name)
	if err != nil {
		return err
	}
	s.Note = strings.TrimSpace(s.Note)
	if len(s.importmap) == 0 {
		s.importmap = map[string]bool{
			"github.com/henrylee2cn/faygo": true,
		}
	}

	var fields = make([]Field, 0, len(s.Fields))
	s.fileParams, s.filesParams = []string{}, []string{}
	for _, field := range s.Fields {
		field.Name = faygo.CamelString(field.Name)
		if field.In == "" {
			continue
		}
		field.isParam = true
		switch field.Type {
		case "*http.Cookie", "http.Cookie":
			s.importmap["net/http"] = true
		case "*multipart.FileHeader", "multipart.FileHeader":
			s.importmap["mime/multipart"] = true
			s.fileParams = append(s.fileParams, field.Name)
		case "[]*multipart.FileHeader", "[]multipart.FileHeader":
			s.importmap["mime/multipart"] = true
			s.filesParams = append(s.filesParams, field.Name)
		}
		fields = append(fields, field)
	}
	s.Fields = fields
	return nil
}

// build struct
func (s *StructHandler) createStruct() string {
	// build struct fields
	var structure string
	structure += fmt.Sprintf("\n/*\n%s %s\n*/\ntype %s struct {", s.Name, s.Note, s.Name)
	for _, field := range s.Fields {
		// if field.Desc != "" {
		// 	structure += fmt.Sprintf("\n  // %s", field.Desc)
		// }
		structure += fmt.Sprintf("\n  %s %s", field.Name, field.Type)
		if !field.isParam {
			if field.OtherTags != "" {
				structure += fmt.Sprintf("`%s`", field.OtherTags)
			}
			continue
		}
		structure += fmt.Sprintf("`param:\"<in:%s>", field.In)
		if field.ParamName != "" {
			structure += fmt.Sprintf("<name:%s>", field.ParamName)
		}
		if field.Required {
			structure += fmt.Sprintf("<required>")
		}
		if field.Nonzero {
			structure += fmt.Sprintf("<nonzero>")
		}
		if field.Len != "" {
			structure += fmt.Sprintf("<len:%s>", field.Len)
		}
		if field.Range != "" {
			structure += fmt.Sprintf("<range:%s>", field.Range)
		}
		if field.Regexp != "" {
			structure += fmt.Sprintf("<regexp:%s>", field.Regexp)
		}
		if field.Maxmb > 0 {
			structure += fmt.Sprintf("<maxmb:%d>", field.Maxmb)
		}
		if field.Err != "" {
			structure += fmt.Sprintf("<err:%s>", field.Err)
		}
		if field.Desc != "" {
			structure += fmt.Sprintf("<desc:%s>", field.Desc)
		}
		structure += fmt.Sprintf("\"")
		if field.OtherTags != "" {
			structure += fmt.Sprintf(" %s", field.OtherTags)
		}
		structure += fmt.Sprintf("`")
		switch field.Type {
		case "*multipart.FileHeader", "multipart.FileHeader":
			structure += fmt.Sprintf("\n  %sUrl string `param:\"-\"`", field.Name)
		case "[]*multipart.FileHeader", "[]multipart.FileHeader":
			structure += fmt.Sprintf("\n  %sUrls []string `param:\"-\"`", field.Name)
		}
	}
	structure += fmt.Sprintf("\n}\n")

	// build methods
	var serve string
	serve += fmt.Sprintf("\n// Serve impletes Handler.\nfunc (%s *%s) Serve(ctx *faygo.Context) error {", s.sign, s.Name)
	if s.ServeContent != "" {
		serve += fmt.Sprintf("\n%s", s.ServeContent)
	} else {
		for i, filename := range s.fileParams {
			var equal = "="
			if i == 0 {
				equal = ":" + equal
			}
			serve += fmt.Sprintf("\n    info, err %s ctx.SaveFile(%q, false)", equal, faygo.SnakeString(filename))
			serve += fmt.Sprintf("\n    if err != nil {\n        return ctx.JSON(412, faygo.Map{\"error\": err.Error()}, true)\n    }")
			serve += fmt.Sprintf("\n    %s.%sUrl = info.Url", s.sign, filename)
		}
		for i, filename := range s.filesParams {
			var equal = "="
			if i == 0 {
				equal = ":" + equal
			}
			serve += fmt.Sprintf("\n    infos, err %s ctx.SaveFiles(%q, false)", equal, faygo.SnakeString(filename))
			serve += fmt.Sprintf("\n    if err != nil {\n        return ctx.JSON(412, faygo.Map{\"error\": err.Error()}, true)\n    }")
			serve += fmt.Sprintf("\n    for _, info := range infos {")
			serve += fmt.Sprintf("\n        %s.%sUrls = append(%s.%sUrls, info.Url)", s.sign, filename, s.sign, filename)
			serve += fmt.Sprintf("\n    }")
		}
		serve += fmt.Sprintf("\n    return ctx.JSON(200, %s, true)", s.sign)
	}
	serve += fmt.Sprintf("\n}\n")

	var doc string
	if s.Note != "" || s.Return != "" {
		doc += fmt.Sprintf("\n// Doc returns the API's note, result or parameters information.\nfunc (%s *%s) Doc() faygo.Doc {", s.sign, s.Name)
		doc += fmt.Sprintf("\n    return faygo.Doc{")
		doc += fmt.Sprintf("\n        Note: %q,", s.Note)
		doc += fmt.Sprintf("\n        Return: %q,", s.Return)
		doc += fmt.Sprintf("\n    }")
		doc += fmt.Sprintf("\n}")
	}

	return structure + serve + doc
}
