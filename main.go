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
	"regexp"
	"strings"
)

//go:embed index.html
var index []byte

//go:embed font.css
var font []byte

//go:embed dark_theme.css
var css []byte

//go:embed favicon.ico
var favicon []byte

var CFG struct {
	Hostname   string
	ListenAddr string
	ListenPort string
	MANPATH    string
	Pandoc     string
}

func cmdout(s string) string {
	ss := strings.Split(s, " ")
	b, e := exec.Command(ss[0], ss[1:]...).Output()
	if e != nil {
		log.Fatal("Fatal: unable to get " + ss[0])
	}
	return strings.TrimSpace(string(b))
}

func init() {
	CFG.MANPATH = cmdout("manpath")
	CFG.Hostname = cmdout("hostname")
	CFG.Pandoc = cmdout("which pandoc")
	CFG.ListenAddr = os.Getenv("ListenAddr")
	CFG.ListenPort = os.Getenv("ListenPort")
	if CFG.ListenPort == "" {
		CFG.ListenPort = "8082"
	}

	css = append(css, font...)
}

func main() {

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprint(len(css)))
		w.WriteHeader(http.StatusOK)
		w.Write(css)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Content-Length", fmt.Sprint(len(favicon)))
		w.WriteHeader(http.StatusOK)
		w.Write(favicon)
	})
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(CFG.ListenAddr+":"+CFG.ListenPort, nil)

}

const htmlHeader = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>{{ title }}</title>
<link rel="stylesheet" type="text/css" href="/style.css">
</head>
<body>
`

func HtmlHeader(body *string, title string) {
	*body = strings.ReplaceAll(htmlHeader, "{{ title }}", strings.ToUpper(title)) + *body + "</body></html>"
	*body = strings.Replace(*body, "<h1>NAME</h1>", "<h1>"+strings.ToUpper(title)+"</h1>", 1)
}

type ManPage struct {
	Section int
	Name    string
	Path    string
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

func ReadFh(path string) (string, error) {
	var buff bytes.Buffer
	fh, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	if strings.HasSuffix(path, ".gz") {
		err = readCompressed(fh, &buff)
		return buff.String(), err
	}
	_, err = buff.ReadFrom(fh)
	return buff.String(), err
}

func pandocConvert(input string) (string, error) {
	cmd := exec.Command(CFG.Pandoc, "--section-divs", "-t", "html4", "-f", "man")
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stdin = strings.NewReader(input)
	b, err := cmd.Output()
	return string(b), err
}

func (m *ManPage) html(w http.ResponseWriter, r *http.Request) error {
	if m.Path == "" {
		return fmt.Errorf("no path")
	}
	var b, fh string
	var err error
	fh, err = ReadFh(m.Path)
	if err != nil {
		return err
	}
	b, err = pandocConvert(fh)
	if err != nil {
		return err
	}
	HtmlHeader(&b, m.Name)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, b)
	return nil
}

func (m *ManPage) FindPath() error {
	s := m.Name
	if m.Section > 0 {
		s += "." + fmt.Sprint(m.Section)
	}
	cmd := exec.Command("man", "-w", s)
	b, e := cmd.Output()
	if e != nil {
		return fmt.Errorf("page not found")
	}
	m.Path = strings.TrimSpace(string(b))
	return nil
}

var manRegexp = []*regexp.Regexp{regexp.MustCompile(`\.[1-9]$`), regexp.MustCompile(`( )?[(][1-9][)]$`)}

func (m *ManPage) ParseName(s string) (err error) {
	s = strings.TrimSpace(s)
	for i, rx := range manRegexp {
		if rx.MatchString(s) {
			m.Section = int((s[len(s)-i-1]) - '0')
			m.Name = strings.TrimSpace(s[:len(s)-i-2])
			return m.FindPath()
		}
	}
	m.Section = 0
	m.Name = s
	return m.FindPath()
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	args := "-l\n-" + strings.Join(r.Form["arg"], "\n-")
	search := strings.ReplaceAll(r.Form["search"][0], "\r", "")
	args += "\n" + search
	cmd := exec.Command("apropos", strings.Split(args, "\n")...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	// cmd.Env = append(cmd.Env, "MANPATH="+CFG.MANPATH)
	b, e := cmd.Output()
	if e != nil {
		http.Error(w, "no results", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strings.ReplaceAll(string(b), "\n", "<br>"))
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
	if err := man.ParseName(q); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes.ReplaceAll(index, []byte("{{ host }}"), []byte(r.Host)))
		return
	}
	fmt.Fprintf(w, "%v", man.html(w, r))
}