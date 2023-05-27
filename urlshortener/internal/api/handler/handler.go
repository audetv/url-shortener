package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/api/vlidator"
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
	r.HandleFunc("/favicon.ico", http.HandlerFunc(r.faviconHandler).ServeHTTP)
	r.HandleFunc("/create", http.HandlerFunc(r.CreateLink).ServeHTTP)
	r.HandleFunc("/redirect", http.HandlerFunc(r.Redirect).ServeHTTP)
	r.HandleFunc("/search", http.HandlerFunc(r.SearchLinks).ServeHTTP)
	r.HandleFunc("/short", http.HandlerFunc(r.SearchShort).ServeHTTP)
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
	Short         shorturl.ShortUrl `json:"short"`
	Search        string            `json:"search"`
	Origin        string            `json:"origin"`
	RedirectCount int               `json:"redirect_count"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
}

type Message struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
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

	// Проверяем, что ссылка правильная
	err := vlidator.ValidLink(l.Origin)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(Message{
			Message: err.Error(),
			Code:    http.StatusUnprocessableEntity,
		})
		return
	}

	ln := link.Link{
		Origin: l.Origin,
	}

	newLink, err := rt.links.CreateLink(r.Context(), ln)
	if err != nil {
		http.Error(w, "error when creating link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(
		Link{
			Short:     newLink.Short,
			Origin:    newLink.Origin,
			CreatedAt: newLink.CreatedAt,
		},
	)
}

func (rt *Router) Redirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s := r.URL.Query().Get("s")
	short := shorturl.Parse(s)

	_, err := rt.links.DoRedirect(r.Context(), *short)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	l, err := rt.links.Read(r.Context(), *short)
	log.Printf("link %v", l)
	if err != nil {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, l.Origin, http.StatusSeeOther)

}

// SearchLinks /search?q=” список всех ссылок, или фильтр ссылок по origin url
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
	w.Header().Set("Content-Type", "application/json")
	first := true
	fmt.Fprintf(w, "[")
	defer fmt.Fprintf(w, "]\r\n")

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
					Short:         l.Short,
					Search:        l.Search,
					Origin:        l.Origin,
					RedirectCount: l.RedirectCount,
					CreatedAt:     l.CreatedAt,
				},
			)
			w.(http.Flusher).Flush()
		}
	}
}

// SearchShort ищет ссылку по коду, возвращает json ссылку или 404 если не найдено.
func (rt *Router) SearchShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s := r.URL.Query().Get("q")

	short := shorturl.Parse(s)

	l, err := rt.links.Read(r.Context(), *short)
	log.Printf("link %v", l)
	if err != nil || l.Short == "" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "[")
	defer fmt.Fprintf(w, "]\r\n")
	_ = enc.Encode(
		Link{
			Short:         l.Short,
			Search:        l.Search,
			Origin:        l.Origin,
			RedirectCount: l.RedirectCount,
			CreatedAt:     l.CreatedAt,
		},
	)
	w.(http.Flusher).Flush()
}

func (rt *Router) faviconHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, "static/favicon.ico")
}
