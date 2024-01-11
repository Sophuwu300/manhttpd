package main

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
)

var userpass map[string]string

func handleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, authOK := r.BasicAuth()
		checksum := sha512.New().Sum([]byte(pass))
		pass = hex.EncodeToString(checksum)
		expectedPass, lookupOK := userpass[user]

		if !authOK || !lookupOK || subtle.ConstantTimeCompare([]byte(expectedPass), []byte(pass)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}