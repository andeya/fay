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
	"errors"
	"fmt"
)

type (
	// Main returns project router codes
	Main struct {
		dir       string // file path or package name
		importmap map[string]bool
		frames    []*Frame
	}
	// Frame faygo app
	Frame struct {
		pkgPrefix string
		name      string
		version   string
		router    *Router
	}
)

// NewMain creates a *Main
func NewMain(dir string) (*Main, error) {
	err := cleanDir(&dir)
	if err != nil {
		return nil, err
	}
	m := &Main{
		dir: dir,
		importmap: map[string]bool{
			"github.com/henrylee2cn/faygo": true,
		},
	}
	return m, nil
}

// AddFrame adds faygo app.
func (m *Main) AddFrame(router *Router, frame string, version ...string) error {
	if frame == "" {
		return errors.New("The frame param must be setted.")
	}
	if router == nil {
		return errors.New("The router param can not be nil.")
	}
	err := router.init()
	if err != nil {
		return err
	}

	var ver string
	if len(version) > 0 {
		ver = version[0]
	}
	var newframe = &Frame{
		name:    frame,
		version: ver,
		router:  router,
	}
	if router.dir != m.dir {
		m.importmap[router.PkgPath()] = true
		newframe.pkgPrefix = router.PkgPrefix()
	} else {
		router.TryMainPkg(m.dir)
	}
	m.frames = append(m.frames, newframe)
	return nil
}

// Output returns main's file.
func (m *Main) Output() error {
	err := writeFile(m.dir, "main.go", m.Create())
	if err != nil {
		return err
	}
	for _, frame := range m.frames {
		err = frame.router.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

// Create returns main's codes.
func (m *Main) Create() string {
	var code string
	code += fmt.Sprintf("\nfunc main() {")

	for _, frame := range m.frames {
		var version string
		if frame.version != "" {
			version = fmt.Sprintf(", %q", frame.version)
		}
		code += fmt.Sprintf("\n    %s%s(faygo.New(%q%s))", frame.pkgPrefix, frame.router.funcname, frame.name, version)
	}
	code += fmt.Sprintf("\n    faygo.Run()")
	code += fmt.Sprintf("\n}\n")

	return fmt.Sprintf("package main\n%s\n%s", importCode(m.importmap), code)
}
