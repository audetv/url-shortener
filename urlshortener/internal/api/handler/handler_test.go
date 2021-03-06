package handler

import (
	"encoding/json"
	"fmt"
	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/db/mem/linkmemstore"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func BuildRouter() *Router {
	linkStore := linkmemstore.NewLinks()
	links := link.NewLinks(linkStore)
	router := NewRouter(links)
	return router
}

func TestRouter_CreateLink(t *testing.T) {
	router := BuildRouter()
	handler := router.AuthMiddleware(http.HandlerFunc(router.CreateLink)).ServeHTTP

	testUrl := "https://test.loc"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(
		http.MethodPost,
		"/create",
		strings.NewReader(fmt.Sprintf(`{"origin":"%v"}`, testUrl)),
	)
	r.SetBasicAuth("admin", "admin")
	handler(w, r)

	if statusCode := w.Code; statusCode != http.StatusCreated {
		t.Errorf("Want status '%d', got '%d'", http.StatusCreated, w.Code)
	}

	dec := json.NewDecoder(w.Body)
	testLink := &Link{}

	if err := dec.Decode(testLink); err != nil {
		t.Errorf("Unable decode response")
	}

	if testLink.Short == "" {
		t.Errorf("Handler returned unexpected empty short link")
	}

	if testLink.Origin != testUrl {
		t.Errorf("Handler returned unexpected orign url: got %v expect %v", testLink.Origin, testUrl)
	}
}

func TestRouter_SearchLinks(t *testing.T) {
	router := BuildRouter()
	handler := router.AuthMiddleware(http.HandlerFunc(router.SearchLinks)).ServeHTTP

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/search?q=test", nil)
	r.SetBasicAuth("admin", "admin")
	handler(w, r)

	if statusCode := w.Code; statusCode != http.StatusOK {
		t.Errorf("Want status '%d', got '%d'", http.StatusOK, w.Code)
	}
}
