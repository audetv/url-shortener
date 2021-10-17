package handler

import (
	"encoding/json"
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
