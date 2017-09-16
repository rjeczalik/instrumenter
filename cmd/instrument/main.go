package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

//go:generate go get github.com/jteeuwen/go-bindata/go-bindata golang.org/x/tools/cmd/goimports
//go:generate go install ../../vendor/golang.org/x/tools/cmd/eg
//go:generate go-bindata -nometadata -o builtin.go -prefix ../../templates/ ../../templates/

const usage = `instrument [package]`

var env = append(os.Environ(), "CGO_ENABLED=0")

func die(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	src, err := gosrc()
	if err != nil {
		die(err)
	}

	var pkg string
	switch args := flag.Args(); len(args) {
	case 0:
		var err error
		if pkg, err = gopkg(src); err != nil {
			die(err)
		}
	case 1:
		pkg = args[0]
	default:
		die(usage)
	}

	if err := instrument(src, pkg); err != nil {
		die(err)
	}
}

func instrument(src, pkg string) error {
	tmp, err := ioutil.TempDir("", "instrument-templates")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	var templates []string
	for _, name := range AssetNames() {
		p, err := Asset(name)
		if err != nil {
			return err
		}
		path := filepath.Join(tmp, name)
		if err := ioutil.WriteFile(path, p, 0644); err != nil {
			return err
		}
		templates = append(templates, path)
	}
	path := filepath.Join(src, pkg)
	for _, template := range templates {
		egArgs := []string{
			"-ignore", "stdlib,github.com/rjeczalik/instrumenter/cmd,github.com/rjeczalik/instrumenter/vendor",
			"-transitive",
			"-w",
			"-t", template,
			pkg,
		}
		if err := run("eg", egArgs...); err != nil {
			return err
		}
		if err := run("goimports", "-w", path); err != nil {
			return err
		}
	}
	return nil
}

func run(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env = env
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func gosrc() (string, error) {
	const sep = string(os.PathSeparator)

	src := os.Getenv("GOPATH")
	if src == "" {
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		src = filepath.Join(u.HomeDir, "go")
	}
	src = filepath.Join(src, "src")
	if !strings.HasSuffix(src, sep) {
		src = src + sep
	}
	return src, nil
}

func gopkg(src string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(wd, src) {
		return "", errors.New("current directory is not inside GOPATH: " + src)
	}
	return strings.TrimPrefix(wd, src), nil
}
