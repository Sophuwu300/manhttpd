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
	CFG.Hostname = os.Getenv("HOSTNAME")
	index = strings.ReplaceAll(index, "{{ hostname }}", CFG.Hostname)
	b, e := exec.Command("which", "mandoc").Output()
	if e != nil || len(b) == 0 {
	    CFG.Mandoc=os.Getenv("MANDOCPATH")
	    if CFG.Mandoc == "" {
    	    log.Fatal("Fatal: no mandoc `apt-get install mandoc`")
    	}
	} else {
	    CFG.Mandoc=strings.TrimSpace(string(b))
	}
	//CFG.Mandoc = "/home/sophie/.local/bin/mandoc"
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
	http.ListenAndServe("0.0.0.0:"+CFG.Addr, nil)
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
var HTMLManName = regexp.MustCompile(`(?:<b>)?([a-zA-Z0-9_.:\-]+)(?:</b>)?\(([0-9][0-9a-z]*)\)`)

type ManPage struct {
	Name    string
	Section string
	Desc    string
	Path    string
}

func (m *ManPage) Where() error {
	var arg = []string{"-w", m.Name}
	if m.Section != "" {
		arg = []string{"-w", "-s" + m.Section, m.Name}
	}
	b, err := exec.Command("man", arg...).Output()
	m.Path = strings.TrimSpace(string(b))
	return err
}
func (m *ManPage) Html() string {
	if m.Where() != nil {
		return fmt.Sprintf("<p>404: Unable to locate page %s</p>", m.Name)
	}
	b, err := exec.Command(CFG.Mandoc, "-Thtml", "-O", "fragment", m.Path).Output()
	if err != nil {
		return fmt.Sprintf("<p>500: server error loading %s</p>", m.Name)
	}
	html := LinkRemover(string(b), "")
	html = HTMLManName.ReplaceAllStringFunc(html, func(s string) string {
		m := HTMLManName.FindStringSubmatch(s)
		return fmt.Sprintf(`<a href="/%s.%s">%s(%s)</a>`, m[1], m[2], m[1], m[2])
	})
	return html
}

var ManDotName = regexp.MustCompile(`^([a-zA-Z0-9_\-]+)(?:\.([0-9a-z]+))?$`)

func NewManPage(s string) (m ManPage) {
	name := ManDotName.FindStringSubmatch(s)
	if len(name) >= 2 {
		m.Name = name[1]
	}
	if len(name) >= 3 {
		m.Section = name[2]
	}
	return
}

var RxWords = regexp.MustCompile(`("[^"]+")|([^ ]+)`).FindAllString
var RxWhatIs = regexp.MustCompile(`([a-zA-Z0-9_\-]+) [(]([0-9a-z]+)[)][\- ]+(.*)`).FindAllStringSubmatch

func searchHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	q := r.Form.Get("q")
	if q == "" || ManDotName.MatchString(q) {
		http.Redirect(w, r, "/"+q, http.StatusFound)
		return
	}
	var args = RxWords("-l "+q, -1)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
		args[i] = strings.TrimPrefix(args[i], `"`)
		args[i] = strings.TrimSuffix(args[i], `"`)
	}
	cmd := exec.Command("apropos", args...)
	b, e := cmd.Output()
	if len(b) < 1 || e != nil {
		WriteHtml(w, r, "Search", fmt.Sprintf("<p>404: no resualts matching %s</p>", q))
		return
	}
	var output string
	for _, line := range RxWhatIs(string(b), -1) { // strings.Split(string(b), "\n") {
		if len(line) == 4 {
			output += fmt.Sprintf(`<p><a href="/%s.%s">%s (%s)</a> - %s</p>%c`, line[1], line[2], line[1], line[2], line[3], 10)
		}
	}
	WriteHtml(w, r, "Search", output)
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
	WriteHtml(w, r, "Index", "")
	return

}