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
	if r.funcname == "" {
		return errors.New("The funcname field must be setted.")
	}
	r.funcname = thinkgo.CamelString(r.funcname)
	if !filepath.IsAbs(r.dir) {
		var err error
		r.dir, err = filepath.Abs(r.dir)
		if err != nil {
			return err
		}
		r.dir = strings.Replace(r.dir, `\`, `/`, -1)
	}
	if len(r.importmap) == 0 {
		r.importmap = map[string]bool{
			"github.com/henrylee2cn/thinkgo": true,
		}
	}
	return nil
}

// API adds handler.
func (r *Router) API(structHandler *StructHandler) error {
	if structHandler == nil {
		return errors.New("The structHandler param can not be nil.")
	}
	err := structHandler.init()
	if err != nil {
		return err
	}
	node := &Node{
		Router:        r,
		structHandler: structHandler,
	}
	pkg := structHandler.PkgPath()
	if pkg != r.PkgPath() {
		r.importmap[pkg] = true
	}
	r.nodes = append(r.nodes, node)
	return nil
}

// CreateFile returns router's file.
func (r *Router) CreateFile() error {
	err := os.MkdirAll(r.dir, 0777)
	if err != nil {
		return err
	}
	filename := path.Join(r.dir, thinkgo.SnakeString(r.funcname)) + ".go"
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	code := r.Create()
	_, err = f.WriteString(code)
	f.Close()
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Run()
	if err != nil {
		return err
	}
	for _, node := range r.nodes {
		if node.structHandler != nil {
			err = node.structHandler.CreateFile()
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

	var imports []string
	for pkg := range r.importmap {
		imports = append(imports, pkg)
	}
	sort.Strings(imports)
	if len(imports) > 0 {
		var pkgs = fmt.Sprintf("\nimport (\n")
		for _, pkg := range imports {
			pkgs += fmt.Sprintf("\n    %q", pkg)
		}
		pkgs += fmt.Sprintf("\n)\n")
		code = pkgs + code
	}

	code = fmt.Sprintf("package %s\n%s", r.PkgName(), code)

	return code
}

// PkgPath returns the package path, e.g `github.com/henrylee2cn/think/test`
func (r *Router) PkgPath() string {
	if r.isMainPkg || r.dir == "" {
		return ""
	}
	dirs := strings.Split(r.dir, "/src/")
	if len(dirs) < 2 {
		return ""
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
		if node.structHandler != nil {
			node.structHandler.TryMainPkg(mainPkgPath)
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
		p := strings.Split(node.structHandler.UrlPath, "/:")[0]
		p = strings.Split(p, "/*")[0]
		ps := strings.Split(p, "/")
		var curNode = root
		var last = len(ps) - 1
	loop:
		for i, p := range ps {
			node.pattern = p
			if i == last {
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

// Node router tree
type Node struct {
	*Router
	pattern       string
	structHandler *StructHandler
	children      []*Node
}

// Create returns struct handler's codes
func (n *Node) Create() string {
	var code string
	code += fmt.Sprintf("\n// Register the route in a tree style.\nfunc Route(frame *thinkgo.Framework) {")
	code += fmt.Sprintf("\nframe.Route(")
	n.create(&code)
	code += fmt.Sprintf("\n)")
	code += fmt.Sprintf("\n}")
	return code
}

func (n *Node) create(code *string) {
	if n.structHandler != nil {
		var pkgPrefix string
		if n.PkgPrefix() != n.structHandler.PkgPrefix() {
			pkgPrefix = n.structHandler.PkgPrefix()
		}
		*code += fmt.Sprintf(
			"\nframe.NewNamedAPI(%q, %q, \"/%s\", &%s%s{}),",
			n.structHandler.RouterName(), n.structHandler.Method, n.pattern, pkgPrefix, n.structHandler.Name,
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
		*code += fmt.Sprintf("\n),")
	}
}
