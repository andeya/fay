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
	"strings"

	"github.com/henrylee2cn/thinkgo"
)

// Router returns project router codes
type Router struct {
	funcname  string
	dir       string // file path or package name
	nodes     []*Node
	isMainPkg bool
	importmap map[string]bool
}

// NewRouter creates a *Router
func NewRouter(funcname, dir string) (*Router, error) {
	r := &Router{
		funcname: funcname,
		dir:      dir,
	}
	err := r.init()
	return r, err
}

func (r *Router) init() error {
	_, err := cleanName(&r.funcname)
	if err != nil {
		return err
	}
	err = cleanDir(&r.dir)
	if err != nil {
		return err
	}
	if len(r.importmap) == 0 {
		r.importmap = map[string]bool{
			"github.com/henrylee2cn/thinkgo": true,
		}
	}
	return nil
}

// AddHandler adds handler.
func (r *Router) AddHandler(handler Handler) error {
	if handler == nil {
		return errors.New("The Handler param can not be nil.")
	}
	_urlPath := handler.GetUrlPath()
	err := handler.init()
	if err != nil {
		return err
	}
	for _, node := range r.nodes {
		if node.urlPath == handler.GetUrlPath() {
			if node.handler != nil || node.static != nil {
				return errors.New("urlPath conflicts: " + _urlPath)
			}
			pkg := handler.PkgPath()
			if pkg != r.PkgPath() {
				r.importmap[pkg] = true
			}
			node.handler = handler
			return nil
		}
	}
	node := &Node{
		Router:  r,
		handler: handler,
		urlPath: handler.GetUrlPath(),
	}
	pkg := handler.PkgPath()
	if pkg != r.PkgPath() {
		r.importmap[pkg] = true
	}
	r.nodes = append(r.nodes, node)
	return nil
}

// AddMiddleware adds midddlerwares.
func (r *Router) AddMiddleware(handlers ...Handler) error {
	if len(handlers) == 0 {
		return nil
	}
loop:
	for _, handler := range handlers {
		if handler == nil {
			return errors.New("The middleware Handler param can not be nil.")
		}
		err := handler.init()
		if err != nil {
			return err
		}
		pkg := handler.PkgPath()
		if pkg != r.PkgPath() {
			r.importmap[pkg] = true
		}
		for _, node := range r.nodes {
			if node.urlPath == handler.GetUrlPath() {
				node.middlewares = append(node.middlewares, handler)
				continue loop
			}
		}
		node := &Node{
			Router:      r,
			urlPath:     handler.GetUrlPath(),
			middlewares: []Handler{handler},
		}
		r.nodes = append(r.nodes, node)
	}
	return nil
}

// AddStatic adds static handler.
func (r *Router) AddStatic(name, urlPath string, root string, nocompressAndNocache ...bool) error {
	_urlPath := urlPath
	err := cleanUrlPath(&urlPath)
	if err != nil {
		return err
	}
	var nocompress, nocache bool
	if len(nocompressAndNocache) > 0 {
		nocompress = nocompressAndNocache[0]
	}
	if len(nocompressAndNocache) > 1 {
		nocache = nocompressAndNocache[1]
	}
	var static = &Static{
		Name:       name,
		UrlPath:    urlPath,
		Root:       root,
		Nocompress: nocompress,
		Nocache:    nocache,
	}
	for _, node := range r.nodes {
		if node.urlPath == urlPath {
			if node.handler != nil || node.static != nil {
				return errors.New("urlPath conflicts: " + _urlPath)
			}
			node.static = static
			return nil
		}
	}
	node := &Node{
		Router:  r,
		urlPath: urlPath,
		static:  static,
	}
	r.nodes = append(r.nodes, node)
	return nil
}

// Output returns router's file.
func (r *Router) Output() error {
	code := r.Create()
	err := writeFile(r.dir, thinkgo.SnakeString(r.funcname)+".go", code)
	if err != nil {
		return err
	}
	for _, node := range r.nodes {
		if node.handler != nil {
			err = node.handler.Output()
			if err != nil {
				return err
			}
		}
		for _, ware := range node.middlewares {
			err = ware.Output()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Create returns router's codes.
func (r *Router) Create() string {
	var code = r.root().Create()
	return fmt.Sprintf("package %s\n%s\n%s", r.PkgName(), importCode(r.importmap), code)
}

// PkgPath returns the package path, e.g `github.com/henrylee2cn/think/test`
func (r *Router) PkgPath() string {
	if r.isMainPkg || r.dir == "" {
		return ""
	}
	dirs := strings.Split(r.dir, "/src/")
	if len(dirs) < 2 {
		thinkgo.Fatalf("You must generate codes in the `src` or its offspring directory!")
	}
	return strings.Join(dirs[1:], "/src/")
}

// TryMainPkg tries to set it as the main package
func (r *Router) TryMainPkg(mainPkgPath string) {
	if r.dir != mainPkgPath {
		return
	}
	r.isMainPkg = true
	for _, node := range r.nodes {
		if node.handler != nil {
			node.handler.TryMainPkg(mainPkgPath)
		}
	}
}

// PkgName returns the package name, e.g `router`
func (r *Router) PkgName() string {
	if r.isMainPkg || r.dir == "" {
		return "main"
	}
	return r.dir[strings.LastIndex(r.dir, "/")+1:]
}

// PkgPrefix returns the package prefix, e.g `router.`
func (r *Router) PkgPrefix() string {
	if r.isMainPkg || r.dir == "" {
		return ""
	}
	return r.dir[strings.LastIndex(r.dir, "/")+1:] + "."
}

func (r *Router) root() *Node {
	var root = &Node{
		Router: r,
	}
	for _, node := range r.nodes {
		var p string
		if node.handler != nil {
			p = strings.Split(node.handler.GetUrlPath(), "/:")[0]
		} else if node.static != nil {
			p = strings.Split(node.static.UrlPath, "/:")[0]
		} else if len(node.middlewares) > 0 {
			p = strings.Split(node.middlewares[0].GetUrlPath(), "/:")[0]
		}
		p = strings.Split(p, "/*")[0]
		var ps = strings.Split(p, "/")
		var curNode = root
		var last = len(ps) - 1
	loop:
		for i, p := range ps {
			node.pattern = p
			if i == last {
				for _, child := range curNode.children {
					if child.pattern == p {
						child.middlewares = append(child.middlewares, node.middlewares...)
						break loop
					}
				}
				curNode.children = append(curNode.children, node)
				break
			}
			for _, child := range curNode.children {
				if child.pattern == p {
					curNode = child
					continue loop
				}
			}
			group := &Node{
				Router:  r,
				pattern: p,
			}
			curNode.children = append(curNode.children, group)
			curNode = group
		}
	}
	return root
}

type (
	// Node router tree
	Node struct {
		*Router
		urlPath     string
		pattern     string
		handler     Handler
		middlewares []Handler
		static      *Static
		children    []*Node
	}
	// Handler interface
	Handler interface {
		Output() error
		TryMainPkg(mainPkgPath string)
		GetUrlPath() string
		PkgPath() string
		PkgPrefix() string
		RouterName() string
		GetName() string
		GetMethod() thinkgo.Methodset
		init() error
	}
)

// Create returns struct handler's codes
func (n *Node) Create() string {
	var code string
	code += fmt.Sprintf("\n// %s register router in a tree style.\nfunc %s(frame *thinkgo.Framework) {", n.funcname, n.funcname)
	code += fmt.Sprintf("\nframe.Route(")
	n.create(&code)
	code += fmt.Sprintf("\n)")
	code += fmt.Sprintf("\n}")
	return code
}

func (n *Node) create(code *string) {
	var use string
	for i, ware := range n.middlewares {
		if i == 0 {
			use += fmt.Sprintf(".Use(")
		}
		var pkgPrefix string
		if n.PkgPrefix() != ware.PkgPrefix() {
			pkgPrefix = ware.PkgPrefix()
		}
		use += fmt.Sprintf("%s%s", pkgPrefix, ware.GetName())
		if i == len(n.middlewares)-1 {
			use += fmt.Sprintf(")")
		} else {
			use += fmt.Sprintf(", ")
		}
	}
	if n.handler != nil {
		var pkgPrefix string
		if n.PkgPrefix() != n.handler.PkgPrefix() {
			pkgPrefix = n.handler.PkgPrefix()
		}
		switch h := n.handler.(type) {
		case *FuncHandler:
			*code += fmt.Sprintf(
				"\nframe.NewNamedAPI(%q, %q, \"/%s\", %s%s)%s,",
				n.handler.RouterName(), n.handler.GetMethod(), n.pattern, pkgPrefix, h.GetName(), use,
			)
		case *StructHandler:
			*code += fmt.Sprintf(
				"\nframe.NewNamedAPI(%q, %q, \"/%s\", &%s%s{})%s,",
				n.handler.RouterName(), n.handler.GetMethod(), n.pattern, pkgPrefix, h.GetName(), use,
			)
		}
		return
	}
	if n.static != nil {
		*code += fmt.Sprintf(
			"\nframe.NewNamedStatic(%q, \"/%s\", %q, %v, %v)%s,",
			n.static.Name, n.pattern, n.static.Root, n.static.Nocompress, n.static.Nocache, use,
		)
		return
	}
	if n.pattern != "" {
		*code += fmt.Sprintf("\nframe.NewGroup(\"/%s\",", n.pattern)
	}
	for _, child := range n.children {
		child.create(code)
	}
	if n.pattern != "" {
		*code += fmt.Sprintf("\n)%s,", use)
	}
}

// Static static router info
type Static struct {
	Name       string
	UrlPath    string
	Root       string
	Nocompress bool
	Nocache    bool
}
