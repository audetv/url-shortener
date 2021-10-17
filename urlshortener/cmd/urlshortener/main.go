package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/audetv/url-shortener/urlshortener/internal/api/handler"
	"github.com/audetv/url-shortener/urlshortener/internal/api/server"
	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/app/starter"
	"github.com/audetv/url-shortener/urlshortener/internal/db/mem/linkmemstore"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ls := linkmemstore.NewLinks()
	app := starter.NewApp(ls)
	links := link.NewLinks(ls)

	h := handler.NewRouter(links)

	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
