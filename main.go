package main

import (
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
var index string

//go:embed dark_theme.css
var css []byte

//go:embed favicon.ico
var favicon []byte

var CFG struct {
	Hostname string
	Addr     string
	Mandoc   string
}

func GetCFG() {
	CFG.Hostname, _ = os.Hostname()
	index = strings.ReplaceAll(index, "{{ hostname }}", CFG.Hostname)
	b, e := exec.Command("which", "mandoc").Output()
	if e != nil || len(b) == 0 {
		log.Fatal("Fatal: no mandoc")
	}
	CFG.Mandoc = strings.TrimSpace(string(b))
	CFG.Addr = os.Getenv("ListenPort")
	if CFG.Addr == "" {
		CFG.Addr = "8082"
	}
}

func main() {
	GetCFG()

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
	http.ListenAndServe(":"+CFG.Addr, nil)

}

func WriteHtml(w http.ResponseWriter, r *http.Request, title, html string) {
	out := strings.ReplaceAll(index, "{{ host }}", r.Host)
	out = strings.ReplaceAll(out, "{{ title }}", title)
	out = strings.ReplaceAll(out, "{{ content }}", html)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, out)
}

var LinkRemover = regexp.MustCompile(`(<a [^>]*>)|(</a>)`).ReplaceAllString
var HTMLManName = regexp.MustCompile(`(<b>)?[a-zA-Z0-9_\-]+(</b>)?\([0-9a-z]+\)`)

type ManPage struct {
	Name    string
	Section string
	Desc    string
}

func (m *ManPage) Path() string {
	arg := "-w"
	if m.Section != "" {
		arg += "s" + m.Section
	}
	b, _ := exec.Command("man", arg, m.Name).Output()
	return strings.TrimSpace(string(b))
}
func (m *ManPage) Html() string {
	b, err := exec.Command(CFG.Mandoc, "-c", "-K", "utf-8", "-T", "html", "-O", "fragment", m.Path()).Output()
	if err != nil {
		return fmt.Sprintf("<p>404: Page %s not found.</p>", m.Name)
	}
	html := LinkRemover(string(b), "")
	return html
}

var ManDotName = regexp.MustCompile(`^([a-zA-Z0-9_\-]+)(?:.([0-9a-z]+))?$`).FindStringSubmatch

func NewManPage(s string) (m ManPage) {
	name := ManDotName(s)
	if len(name) >= 2 {
		m.Name = name[1]
	}
	if len(name) >= 3 {
		m.Section = name[2]
	}
	return
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
		// http.Redirect(w, r, "/", http.StatusFound)
		man := NewManPage(r.URL.Path[1:])
		WriteHtml(w, r, man.Name, man.Html())
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
	var man ManPage
	man.Name = r.Form.Get("man")
	if man.Name == "" {
		WriteHtml(w, r, "Index", "")
		return
	}

}