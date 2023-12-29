package main

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	http.HandleFunc("/cgi-bin/man/", handleWIS)
	log.Println("Starting server on http://localhost:1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

func handleWIS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	exe := filepath.Base(r.URL.Path)
	if !(exe == "man2html" || exe == "mansearch" || exe == "mansec" || exe == "manwhatis") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	q := ""
	if strings.Contains(r.URL.String(), "?") {
		q = strings.SplitN(r.URL.String(), "?", 2)[1]
	}
	cmd := exec.Command("/usr/lib/cgi-bin/man/" + exe)
	var buff bytes.Buffer
	cmd.Env = append(cmd.Env, "QUERY_STRING="+q, "REQUEST_METHOD="+r.Method, "SERVER_NAME=localhost:1234")
	// cmd.Env = append(cmd.Env, "MANPATH=/usr/man:/usr/share/man:/usr/local/man:/usr/local/share/man:/usr/X11R6/man:/opt/man:/snap/man")
	cmd.Stdout = &buff
	err := cmd.Run()
	if err != nil {
		http.Redirect(w, r, "/cgi-bin/man/man2html", http.StatusTemporaryRedirect)
		return
	}
	buff.WriteTo(w)
}