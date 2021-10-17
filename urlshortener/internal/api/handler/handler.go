package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

type Router struct {
	*http.ServeMux
	links *link.Links
}

func NewRouter(links *link.Links) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		links:    links,
	}
	r.HandleFunc("/create", r.AuthMiddleware(http.HandlerFunc(r.CreateLink)).ServeHTTP)
	r.HandleFunc("/search", r.AuthMiddleware(http.HandlerFunc(r.SearchLinks)).ServeHTTP)
	return r
}

func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

type Link struct {
	Short  shorturl.ShortUrl `json:"short"`
	Origin string            `json:"origin"`
}

func (rt *Router) CreateLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	l := Link{}

	if err := dec.Decode(&l); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if l.Origin == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// TODO Сделать проверку, что урл валидный

	ln := link.Link{
		Origin: l.Origin,
	}

	newLink, err := rt.links.CreateLink(r.Context(), ln)
	if err != nil {
		http.Error(w, "error when creating link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(Link{
		Short:  newLink.Short,
		Origin: newLink.Origin,
	})
}

// SearchLinks /search?q='' список всех ссылок, или фильтр ссылок по origin url
func (rt *Router) SearchLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")

	ch, err := rt.links.SearchLinks(r.Context(), q)
	if err != nil {
		http.Error(w, "error when searching", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	first := true
	fmt.Printf("[")
	defer fmt.Printf("]\r\n")

	for {
		select {
		case <-r.Context().Done():
			return
		case l, ok := <-ch:
			if !ok {
				return
			}
			if first {
				first = false
			} else {
				fmt.Fprintf(w, ",")
			}
			_ = enc.Encode(
				Link{
					Short:  l.Short,
					Origin: l.Origin,
				},
			)
			w.(http.Flusher).Flush()
		}
	}
}
