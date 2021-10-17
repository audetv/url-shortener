package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
)

type Server struct {
	srv   http.Server
	links *link.Links
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Start(links *link.Links) {
	s.links = links
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			log.Printf("serve error %v:", err)
		}
	}()
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_ = s.srv.Shutdown(ctx)
	cancel()
}
