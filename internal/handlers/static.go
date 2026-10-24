package handlers

import (
	_ "embed"
	"net/http"
)

//go:embed res/index.html
var indexHTML []byte

func BasicPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(indexHTML)
}
