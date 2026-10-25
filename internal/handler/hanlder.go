package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/DarkSeriusCode/url-shortener/internal/service"
)

var (
	indexHTML = template.Must(template.ParseFiles("ui/index.html"))
	staticFiles = http.FileServer(http.Dir("ui/static/"))
)

type AppHandler struct {
	shortener *service.Shortener
	mux       *http.ServeMux
}

func (h *AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func NewHandler(shortener *service.Shortener) *AppHandler {
	mux := http.NewServeMux()
	h := &AppHandler{
		shortener: shortener,
		mux:       mux,
	}

	mux.HandleFunc("/", h.mainPage)
	mux.HandleFunc("/shorten", h.shorten)
	mux.HandleFunc("/{id}", h.resolve)
	mux.Handle("/static/", http.StripPrefix("/static", staticFiles))

	return h
}

func (h *AppHandler) mainPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err := indexHTML.Execute(w, nil); err != nil {
		slog.Error(fmt.Sprintf("/ Error while executing HTML template %s", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	slog.Info("/")
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
	urlID, err := h.shortener.Shorten(*originURL)
	if err != nil {
		http.Error(w, "Databse error", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("/shorten Database error: %s", err))
		return
	}

	resp := URLShortenResp{ID: urlID}
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

func (h *AppHandler) resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	shortenedURL, err := h.shortener.Resolve(id)
	if err != nil {
		slog.Error(fmt.Sprintf("Not such URL with ID %s", id), slog.String("id", id))
		http.Error(w, "No such URL ID", http.StatusNotFound)
		return
	}
	slog.Info("/resolve Redirecting",
		slog.String("url", shortenedURL.OriginURL.String()),
		slog.String("id", shortenedURL.UrlID),
	)
	http.Redirect(w, r, shortenedURL.OriginURL.String(), http.StatusPermanentRedirect)
}
