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

package model

import (
	"strings"

	"github.com/henrylee2cn/think/generator"
	"github.com/henrylee2cn/thinkgo"
)

// SimplePro output project files.
func SimplePro(projectDir string, appname string, appVersion ...string) {
	initDir(projectDir)

	router, err := generator.NewRouter("Route", projectDir+"router")
	if err != nil {
		thinkgo.Fatalf("[think] Create project fail:%v", err)
	}
	router.AddHandler(indexHandler)
	router.AddHandler(testHandler)
	router.AddMiddleware(tokenWare)

	project, err := generator.NewMain(projectDir)
	if err != nil {
		thinkgo.Fatalf("[think] Create project fail:%v", err)
	}
	project.AddFrame(router, appname, appVersion...)
	err = project.Output()
	if err != nil {
		thinkgo.Fatalf("[think] Create project fail:%v", err)
	}
}

func initDir(dir string) {
	indexHandler.Dir = strings.Replace(indexHandler.Dir, "<<DIR>>", dir, -1)
	testHandler.Dir = strings.Replace(testHandler.Dir, "<<DIR>>", dir, -1)
	tokenWare.Dir = strings.Replace(tokenWare.Dir, "<<DIR>>", dir, -1)
}

var indexHandler = &generator.FuncHandler{
	Dir:     "<<DIR>>handler",
	Name:    "Index",
	UrlPath: "/",
	Method:  "GET",
	ServeContent: `return ctx.Render(200, thinkgo.JoinStatic("index.html"), thinkgo.Map{
            "TITLE":   "thinkgo",
            "VERSION": thinkgo.VERSION,
            "CONTENT": "Welcome To Thinkgo",
            "AUTHOR":  "HenryLee",
        })`,
}

var testHandler = &generator.StructHandler{
	Dir:     "<<DIR>>handler",
	Name:    "Test",
	UrlPath: "/test",
	Method:  "POST",
	Fields: []generator.Field{
		{
			Type: "string",
			Name: "Token",
			In:   "query",
		},
		{
			Type:     "string",
			Name:     "Name",
			In:       "formData",
			Required: true,
			Len:      "1:10",
			Desc:     "your name (1~10 words)",
		},
		{
			Type:  "uint8",
			Name:  "Age",
			In:    "formData",
			Range: "1:100",
			Desc:  "your age (1~100)",
		},
		{
			Type: "*multipart.FileHeader",
			Name: "Avatar",
			In:   "formData",
			Desc: "your avatar",
		},
	},
	Note:   "test struct handler",
	Return: "// JSON\n{}",
}

var tokenWare = &generator.FuncHandler{
	Dir:     "<<DIR>>middleware",
	Name:    "Token",
	UrlPath: "/test",
	ServeContent: `ctx.Log().Debugf("[ware] token:%q", ctx.QueryParam("token"))
    return nil`,
}
