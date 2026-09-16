package handler

import "net/http"

func NewHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", mainPage)

	return mux
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Hello world"))
}
