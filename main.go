package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/cgi-bin/man/", handleWIS)
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		b, _ := os.ReadFile("style.css")
		w.Write(b)
		b, _ = os.ReadFile("font.css")
		w.Write(b)
	})
	http.Handle("/", http.RedirectHandler("/cgi-bin/man/man2html", http.StatusTemporaryRedirect))
	http.ListenAndServe("0.0.0.0:1234", nil)
}

func handleWIS(w http.ResponseWriter, r *http.Request) {
	exe := filepath.Base(r.URL.Path)
	if !(exe == "man2html" || exe == "mansearch" || exe == "mansec" || exe == "manwhatis") {
		http.Redirect(w, r, "/cgi-bin/man/man2html", http.StatusTemporaryRedirect)
		return
	}
	q := ""
	if strings.Contains(r.URL.String(), "?") {
		q = strings.SplitN(r.URL.String(), "?", 2)[1]
	}
	cmd := exec.Command("/usr/lib/cgi-bin/man/" + exe)
	var buff bytes.Buffer
	cmd.Env = append(cmd.Env, "QUERY_STRING="+q, "REQUEST_METHOD="+r.Method, "SERVER_NAME=localhost:1234")
	cmd.Env = append(cmd.Env, "MANPATH=/usr/man:/usr/share/man:/usr/local/man:/usr/local/share/man:/usr/X11R6/man:/opt/man:/snap/man")
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		http.Redirect(w, r, "/cgi-bin/man/man2html", http.StatusTemporaryRedirect)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	page := buff.String()
	page = page[strings.Index(page, "<!"):]
	i := strings.Index(page, "</HEAD>")
	fmt.Fprint(w, page[:i], `<link rel="stylesheet" href="/style.css">`, page[i:])
}