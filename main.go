package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/mholt/archiver/v4"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var m2h = "/home/sophuwu/Documents/project/manpages/build/mh.2"

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
	MANPATH    string
}

func cmdout(s string) string {
	ss := strings.Split(s, " ")
	b, e := exec.Command("/usr/bin/"+ss[0], ss[1:]...).Output()
	if e != nil {
		log.Fatal("Fatal: unable to get " + ss[0])
	}
	return strings.TrimSpace(string(b))
}

func init() {
	CFG.MANPATH = cmdout("manpath -g")
	CFG.Hostname = cmdout("hostname")
	CFG.ListenAddr = os.Getenv("ListenAddr")
	CFG.ListenPort = os.Getenv("ListenPort")
	if CFG.ListenPort == "" {
		CFG.ListenPort = "8080"
	}

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
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(CFG.ListenAddr+":"+CFG.ListenPort, nil)

}

type ManPage struct {
	Section string
	Name    string
	Path    string
	WhatIs  string
}

func readCompressed(fh *os.File, buff *bytes.Buffer) error {
	decompressor, err := archiver.Gz{}.OpenReader(fh)
	if err != nil {
		return err
	}
	defer decompressor.Close()
	_, err = buff.ReadFrom(decompressor)
	return err
}

func ReadFh(buff *bytes.Buffer, path string) error {
	fh, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fh.Close()
	if strings.HasSuffix(path, ".gz") {
		return readCompressed(fh, buff)
	}
	_, err = buff.ReadFrom(fh)
	return err
}

func (m *ManPage) html(w http.ResponseWriter) error {
	if m.Path == "" {
		return fmt.Errorf("no path")
	}
	var buff bytes.Buffer
	err := ReadFh(&buff, m.Path)
	if err != nil {
		return err
	}
	cmd := exec.Command(m2h, "-H", CFG.Hostname+":"+CFG.ListenPort, "-M", "/", "-")
	cmd.Stdin = &buff
	b, e := cmd.Output()
	if e != nil {
		return fmt.Errorf("page not found")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	return nil
}

func (m *ManPage) FindPath() error {
	s := m.Name
	if m.Section != "" {
		s = m.Section + "." + s
	}
	cmd := exec.Command("man", "-w", s)
	cmd.Env = append(os.Environ(), "MANPATH="+CFG.MANPATH)
	b, e := cmd.Output()
	if e != nil {
		return fmt.Errorf("page not found")
	}
	m.Path = strings.TrimSpace(string(b))
	return nil
}

func (m *ManPage) FindHumanInput(s string) error {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "(", "")
	s = strings.ReplaceAll(s, ")", "")
	if strings.Contains(s, " ") {
		arr := strings.SplitN(s, " ", 2)
		m.Section, m.Name = arr[1], arr[0]
	} else if strings.Contains(s, ".") {
		arr := strings.SplitN(s, ".", 2)
		m.Section, m.Name = arr[0], arr[1]
	} else {
		m.Name = s
	}
	return m.FindPath()
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	args := "-l\n-" + strings.Join(r.Form["arg"], "\n-")
	search := strings.ReplaceAll(r.Form["search"][0], "\r", "")
	args += "\n" + search
	cmd := exec.Command("apropos", strings.Split(args, "\n")...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "MANPATH="+CFG.MANPATH)
	b, e := cmd.Output()
	if e != nil {
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if r.Method == "POST" {
		searchHandler(w, r)
		return
	}
	if !strings.HasPrefix(r.URL.RawQuery, "man=") && r.URL.RawQuery != "" {
		r.URL.RawQuery = "man=" + r.URL.RawQuery
	}
	_ = r.ParseForm()
	q := r.Form.Get("man")
	var man ManPage
	if err := man.FindHumanInput(q); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(index)
		return
	}
	fmt.Fprintf(w, "%v", man.html(w))
}