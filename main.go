package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//go:embed index.html
var index []byte

//go:embed dark_theme.css
var css []byte

//go:embed favicon.ico
var favicon []byte

var CFG struct {
	Hostname   string
	ListenAddr string
	ListenPort string
	Mandoc     string
}

func init() {
	CFG.Hostname, _ = os.Hostname()
	b, e := exec.Command("which", "mandoc").Output()
	if e != nil || len(b) == 0 {
		log.Fatal("Fatal: no mandoc")
	}
	CFG.Mandoc = strings.TrimSpace(string(b))
	CFG.ListenAddr = os.Getenv("ListenAddr")
	CFG.ListenPort = os.Getenv("ListenPort")
	if CFG.ListenPort == "" {
		CFG.ListenPort = "8082"
	}
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

type ManPage struct {
	Section int
	Name    string
	Path    string
}

func (m *ManPage) html(w http.ResponseWriter, r *http.Request) error {
	if m.Path == "" {
		return fmt.Errorf("no path")
	}
	b, err := exec.Command(CFG.Mandoc, "-Thtml", m.Path).Output()
	if err != nil {
		return err
	}
	s := string(b)
	// HtmlHeader(&s, m.Name)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, s)
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