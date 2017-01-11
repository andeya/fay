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
)

type (
	// Main returns project router codes
	Main struct {
		dir       string // file path or package name
		importmap map[string]bool
		frames    []*Frame
	}
	// Frame thinkgo app
	Frame struct {
		pkgPrefix string
		name      string
		version   string
		router    *Router
	}
)

// NewMain creates a *Main
func NewMain(dir string) (*Main, error) {
	if !filepath.IsAbs(dir) {
		var err error
		dir, err = filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		dir = strings.Replace(dir, `\`, `/`, -1)
	}
	m := &Main{
		dir: dir,
		importmap: map[string]bool{
			"github.com/henrylee2cn/thinkgo": true,
		},
	}
	return m, nil
}

// AddFrame adds thinkgo app.
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

// CreateFile returns main's file.
func (m *Main) CreateFile() error {
	err := os.MkdirAll(m.dir, 0777)
	if err != nil {
		return err
	}
	filename := path.Join(m.dir, "main.go")
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	code := m.Create()
	_, err = f.WriteString(code)
	f.Close()
	cmd := exec.Command("gofmt", "-w", filename)
	err = cmd.Run()
	if err != nil {
		return err
	}
	for _, frame := range m.frames {
		err = frame.router.CreateFile()
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
		code += fmt.Sprintf("\n    %s%s(thinkgo.New(%q%s))", frame.pkgPrefix, frame.router.funcname, frame.name, version)
	}
	code += fmt.Sprintf("\n    thinkgo.Run()")
	code += fmt.Sprintf("\n}\n")

	var imports []string
	for pkg := range m.importmap {
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

	code = fmt.Sprintf("package main\n%s", code)

	return code
}
