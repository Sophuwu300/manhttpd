package main

import (
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"os/exec"
	"path/filepath"
	"strings"
)

func Cmd(s string) []byte {
	arg := strings.Split(s, " ")
	cmd := exec.Command(arg[0], arg[1:]...)
	out, _ := cmd.CombinedOutput()
	return out
}

func Err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fs := memoryfs.New()
	var err error
	var file string

	file = filepath.Join(fs.FSTempDir(), "test.txt")

	err = vfs.WriteFile(fs, file, Cmd("man --pager=cat man"), 0644)
	Err(err)

	var b []byte
	b, err = vfs.ReadFile(fs, file)
	Err(err)
	println(string(b))
}