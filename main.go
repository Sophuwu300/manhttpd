package main

import (
	"bytes"
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
	http.Handle("/", http.RedirectHandler("/cgi-bin/man/man2html", http.StatusTemporaryRedirect))
	http.ListenAndServe("0.0.0.0:3234", nil)
}

func getARG(exe, pth string) string {
	i := strings.Index(pth, exe)
	if i == -1 {
		return ""
	}
	return pth[i+len(exe)-1:]
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

	fmt.Println(r.Method, r.URL.Path, r.URL.Query().Encode())

	exe := getEXE(r.URL.Path)
	if exe == "" {
		http.RedirectHandler("/cgi-bin/man/man2html", 404)
	}

	q := "QUERY_STRING=" + r.URL.RawQuery

	var opt []string
	var rx string = fmt.Sprint(r.URL.Query()["query"])
	if strings.IndexAny(rx, `^*&|;`) > 0 {
		rx = fmt.Sprintf("QUERY_STRING= -r'%s' ", rx)
		opt = append(opt, rx)
		q = rx + q
		exe = "mansearch"
	}

	cmd := exec.Command("/usr/lib/cgi-bin/man/"+exe, opt...)

	var buff bytes.Buffer
	cmd.Stdout = &buff

	cmd.Env = append(cmd.Env, q, "REQUEST_METHOD="+r.Method, "SERVER_NAME=localhost:1234")
	cmd.Env = append(cmd.Env, "MANPATH=/usr/man:/usr/share/man:/usr/local/man:/usr/local/share/man:/usr/X11R6/man:/opt/man:/snap/man")

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	page := buff.String()
	page = page[strings.Index(page, "<!"):]
	i := strings.Index(page, "</HEAD>")
	if i == -1 {
		i = strings.Index(page, "</head>")
	}
	fmt.Fprintln(w, page[:i])
	fmt.Fprintln(w, `<link rel="stylesheet" href="/style.css">`)
	fmt.Fprintln(w, page[i:])
}