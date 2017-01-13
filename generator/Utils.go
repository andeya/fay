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
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/henrylee2cn/thinkgo"
)

// Output returns file.
func Output(fullfilename string, content string) error {
	fullfilename, err := filepath.Abs(fullfilename)
	if err != nil {
		return err
	}
	fullfilename = strings.Replace(fullfilename, `\`, `/`, -1)
	dir, shorname := filepath.Split(fullfilename)
	return writeFile(dir, shorname, content)
}

func writeFile(dir, shortname, code string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	filename := path.Join(dir, shortname)
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = f.WriteString(code)
	f.Close()
	thinkgo.Printf("[think] Created %s", filename)
	if filepath.Ext(shortname) == ".go" {
		cmd := exec.Command("gofmt", "-w", filename)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}

func cleanDir(dir *string) error {
	if !filepath.IsAbs(*dir) {
		var err error
		*dir, err = filepath.Abs(*dir)
		if err != nil {
			return err
		}
	}
	*dir = strings.Replace(*dir, `\`, `/`, -1)
	return nil
}

func cleanUrlPath(urlPath *string) error {
	*urlPath = strings.TrimSpace(*urlPath)
	*urlPath = strings.Trim(*urlPath, "/")
	return nil
}

func cleanName(name *string) (firstChar string, err error) {
	_name := *name
	*name = strings.TrimSpace(*name)
	*name = thinkgo.CamelString(*name)
	if *name == "" {
		err = errors.New("The type (or func) name \"" + _name + "\" is incorrect.")
		return
	}
	signRune := []rune(*name)[0]
	if signRune < 'A' || signRune > 'Z' {
		err = errors.New("The type (or func) name \"" + _name + "\" is incorrect.")
		return
	}
	firstChar = strings.ToLower(string(signRune))
	return
}

func importCode(importmap map[string]bool) string {
	var imports []string
	for pkg := range importmap {
		imports = append(imports, pkg)
	}
	sort.Strings(imports)
	var pkgs string
	if len(imports) > 0 {
		pkgs = fmt.Sprintf("\nimport (\n")
		for _, pkg := range imports {
			pkgs += fmt.Sprintf("\n    %q", pkg)
		}
		pkgs += fmt.Sprintf("\n)\n")
	}
	return pkgs
}
