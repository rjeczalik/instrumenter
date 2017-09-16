package main

//go:generate go get github.com/jteeuwen/go-bindata/go-bindata golang.org/x/tools/cmd/goimports
//go:generate go install ../../vendor/golang.org/x/tools/cmd/eg
//go:generate go-bindata -nometadata -o builtin.go -prefix ../../templates/ ../../templates/

func main() {
}
