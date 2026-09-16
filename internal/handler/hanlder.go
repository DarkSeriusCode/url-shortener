package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/DarkSeriusCode/url-shortener/internal/service"
)

type AppHandler struct {
	shortener *service.Shortener
	mux *http.ServeMux
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func NewHandler(shortener *service.Shortener) *AppHandler {
	mux := http.NewServeMux()
	h := &AppHandler{
		shortener: shortener,
		mux: mux,
	}

	mux.HandleFunc("/", h.mainPage)
	mux.HandleFunc("/shorten", h.shorten)

	return h
}

func (h *AppHandler) mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Write([]byte("Hello world"))
}

func (h *AppHandler) shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		slog.Warn("/shorten Unsupported method", slog.String("method", r.Method))
		return
	}
	reqJSON, err := ReadJSONFromRequest[URLShortenReq](r)
	if err != nil {
		slog.Error(fmt.Sprintf("/shorten %s", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validating
	originURL, err := url.Parse(reqJSON.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		slog.Error("/shorten Invalid URL to short", slog.String("url", reqJSON.URL))
		return
	}
	urlID := h.shortener.Shorten(*originURL)

	resp := URLShortenResp{ ID: urlID }
	bytes, err := json.Marshal(resp)
	if err != nil {
		slog.Error(fmt.Sprintf("/shorten Unable to response %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(bytes)
	slog.Info("/shorten",
		slog.String("url", originURL.String()),
		slog.String("id", urlID),
	)
}
