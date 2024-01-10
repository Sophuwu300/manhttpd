package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

//go:embed font.css
var font []byte

//go:embed dark_theme.css
var css []byte

func init() {
	css = append(css, font...)
}

func main() {
	http.HandleFunc("/cgi-bin/man/", handleWIS)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Content-Length", fmt.Sprint(len(css)))
		w.WriteHeader(http.StatusOK)
		w.Write(css)
	})
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/", handleTMP)
	http.ListenAndServe("0.0.0.0:3234", nil)
}

func getARG(exe, pth string) string {
	i := strings.Index(pth, exe)
	if i == -1 || len(pth) <= i+len(exe) {
		return ""
	}
	return pth[i+len(exe):]
}

func getEXE(s string) string {
	mandex := [4]string{"man2html", "manwhatis", "mansearch", "mansec"}
	for k := range mandex {
		if strings.Contains(s, mandex[k]) {
			return mandex[k]
		}
	}
	return ""
}

func handleWIS(w http.ResponseWriter, r *http.Request) {

	var opt []string

	exe := getEXE(r.URL.Path)
	if exe == "" {
		http.RedirectHandler("/cgi-bin/man/man2html", 404)
	}
	if s := getARG(exe, r.URL.Path); s != "" {
		opt = append(opt, s)
	}

	q := "QUERY_STRING=" + r.URL.RawQuery

	r.ParseForm()
	rx := r.Form.Get("query")
	if strings.ContainsAny(rx, `*&|`) {
		// rx = fmt.Sprintf(`%s`, rx)
		opt = append(opt, rx)
		exe = "mansearch"
	}

	cmd := exec.Command("/usr/lib/cgi-bin/man/"+exe, opt...)

	cmd.Env = append(cmd.Env, q)
	// cmd.Env = append(cmd.Env, "MANPATH=/usr/man:/usr/share/man:/usr/local/man:/usr/local/share/man:/usr/X11R6/man:/opt/man:/snap/man")

	b, _ := cmd.CombinedOutput()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	page := string(b)
	page = page[strings.Index(page, "<!"):]
	i := strings.Index(page, "</HEAD>")
	if i == -1 {
		i = strings.Index(page, "</head>")
	}
	fmt.Fprintln(w, page[:i])
	fmt.Fprintln(w, `<link rel="stylesheet" href="/style.css">`)
	fmt.Fprintln(w, page[i:])
}