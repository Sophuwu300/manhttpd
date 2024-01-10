package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

//go:embed index.html
var index []byte

//go:embed font.css
var font []byte

//go:embed dark_theme.css
var css []byte

var CFG struct {
	Hostname   string
	ListenAddr string
	ListenPort string
	MANPATH    []string
}

func init() {
	CFG.Hostname, _ = os.Hostname()
	if CFG.Hostname == "" {
		os.Getenv("HOSTNAME")
	}
	CFG.ListenAddr = os.Getenv("ListenAddr")
	CFG.ListenPort = os.Getenv("ListenPort")
	if CFG.ListenAddr == "" || CFG.ListenPort == "" || CFG.Hostname == "" {
		log.Fatal("ListenAddr, ListenPort and Hostname must be set")
	}
	b, err := exec.Command("/usr/bin/manpath", "-g").Output()
	if err != nil {
		log.Fatal("Fatal: unable to get manpath")
	}
	CFG.MANPATH = strings.Split(string(b), ":")

	css = append(css, font...)
	index = bytes.ReplaceAll(index, []byte("{{ hostname }}"), []byte(CFG.Hostname))
	index = bytes.ReplaceAll(index, []byte("{{ port }}"), []byte(CFG.ListenPort))
}

func main() {
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprint(len(css)))
		w.WriteHeader(http.StatusOK)
		w.Write(css)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", handler)
	http.ListenAndServe(CFG.ListenAddr+":"+CFG.ListenPort, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	_ = r.ParseForm()

	man := r.Form.Get("man")
	if man == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(index)
		return
	}

	// cmd := exec.Command(man)

	// b, _ := cmd.CombinedOutput()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}