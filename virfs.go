package main

/*
import (
	"fmt"
	"github.com/mandelsoft/vfs/pkg/memoryfs"
	"github.com/mandelsoft/vfs/pkg/vfs"
	"os/exec"
	"strings"
)

func Cmd(s string) []byte {
	arg := strings.Split(s, " ")
	cmd := exec.Command(arg[0], arg[1:]...)
	out, err := cmd.CombinedOutput()
	Err(err)
	return out
}

func Err(err error) {
	if err != nil {
		panic(err)
	}
}

func printDir(dir string) {
	fi, err := vfs.ReadDir(fs, dir)
	Err(err)
	fmt.Println("reading dir: ", dir)
	for _, f := range fi {
		fmt.Println(f.Name(), f.IsDir(), f.Size())
	}
}

var fs vfs.FileSystem

func main() {
	fs = memoryfs.New()
	wd, err := fs.Getwd()
	Err(err)

	printDir(wd)

}

*/