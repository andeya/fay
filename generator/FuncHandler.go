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

	"github.com/henrylee2cn/thinkgo"
)

type (
	// FuncHandler function handler
	FuncHandler struct {
		Dir          string // file path or package name
		UrlPath      string // URL's path
		Name         string // (required) struct name
		Note         string // note for API
		ServeContent string // main logic
		Return       string // response content demo
		Method       thinkgo.Methodset
		fileParams   []string
		filesParams  []string
		importmap    map[string]bool
		sign         string
		isMainPkg    bool
	}
)

// Output creates struct handler file.
func (s *FuncHandler) Output() error {
	code, err := s.Create()
	if err != nil {
		return err
	}
	return writeFile(s.Dir, thinkgo.SnakeString(s.Name)+".go", code)
}

// Create returns struct handler's codes
func (s *FuncHandler) Create() (code string, err error) {
	// initialize
	err = s.init()
	if err != nil {
		return
	}
	// build codes
	code = fmt.Sprintf("package %s\n%s\n%s", s.PkgName(), importCode(s.importmap), s.createFunc())
	return code, nil
}

// GetUrlPath returns router node's url path.
func (s *FuncHandler) GetUrlPath() string {
	return s.UrlPath
}

// GetMethod returns request method.
func (s *FuncHandler) GetMethod() thinkgo.Methodset {
	return s.Method
}

// GetName returns handler type name.
func (s *FuncHandler) GetName() string {
	return s.Name
}

// RouterName returns router node's name
func (s *FuncHandler) RouterName() string {
	if len(s.Note) > 0 {
		return strings.Split(s.Note, "\n")[0]
	}
	return s.Name
}

// TryMainPkg tries to set it as the main package
func (s *FuncHandler) TryMainPkg(mainPkgPath string) {
	if s.Dir != mainPkgPath {
		return
	}
	s.isMainPkg = true
}

// PkgPath returns the package path, e.g `github.com/henrylee2cn/think/test`
func (s *FuncHandler) PkgPath() string {
	if s.isMainPkg || s.Dir == "" {
		return ""
	}
	dirs := strings.Split(s.Dir, "/src/")
	if len(dirs) < 2 {
		thinkgo.Fatalf("You must generate codes in the `src` or its offspring directory!")
	}
	return strings.Join(dirs[1:], "/src/")
}

// PkgName returns the package name, e.g `handler`
func (s *FuncHandler) PkgName() string {
	if s.isMainPkg || s.Dir == "" {
		return "main"
	}
	return s.Dir[strings.LastIndex(s.Dir, "/")+1:]
}

// PkgPrefix returns the package name, e.g `handler.`
func (s *FuncHandler) PkgPrefix() string {
	if s.isMainPkg || s.Dir == "" {
		return ""
	}
	return s.Dir[strings.LastIndex(s.Dir, "/")+1:] + "."
}

// initialize
func (s *FuncHandler) init() error {
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
			"github.com/henrylee2cn/thinkgo": true,
		}
	}
	return nil
}

// build FuncHandler
func (s *FuncHandler) createFunc() string {
	var function string
	var hasdoc = s.Note != "" || s.Return != ""
	if hasdoc {
		function += fmt.Sprintf("\n/*\n%s %s\n*/\nvar %s = thinkgo.WrapDoc(", s.Name, s.Note, s.Name)
		function += fmt.Sprintf("\nthinkgo.HandlerFunc(func(ctx *thinkgo.Context) error {")
	} else {
		function += fmt.Sprintf("\n/*\n%s %s\n*/\nvar %s = thinkgo.HandlerFunc(func(ctx *thinkgo.Context) error {", s.Name, s.Note, s.Name)
	}
	if s.ServeContent == "" {
		function += fmt.Sprintf("\nthinkgo.Debug(\"Calling FuncHandler - %s...\")", s.Name)
		function += fmt.Sprintf("\nreturn nil")
	} else {
		function += fmt.Sprintf("\n%s", s.ServeContent)
	}
	if hasdoc {
		function += fmt.Sprintf("\n}),")
		function += fmt.Sprintf("\n%q,", s.Note)
		function += fmt.Sprintf("\n%q,", s.Return)
		function += fmt.Sprintf("\n)\n")
	} else {
		function += fmt.Sprintf("\n})\n")
	}
	return function
}
