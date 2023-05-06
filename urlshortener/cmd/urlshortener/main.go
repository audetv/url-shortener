package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/secret"

	"github.com/audetv/url-shortener/urlshortener/internal/db/mem/linkmemstore"

	"github.com/audetv/url-shortener/urlshortener/internal/db/sql/pgstore"

	"github.com/audetv/url-shortener/urlshortener/internal/api/handler"
	"github.com/audetv/url-shortener/urlshortener/internal/api/server"
	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/app/starter"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	// output current time zone
	tnow := time.Now()
	tz, _ := tnow.Zone()
	log.Printf("Local time zone %s. Service started at %s", tz,
		tnow.Format("2006-01-02T15:04:05.000 MST"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	var lst link.LinkStoreInterface
	stl := os.Getenv("URL_SHORTENER_STORE")

	username := os.Getenv("POSTGRES_USER")
	passwordBytes := secret.Read(os.Getenv("POSTGRES_PASSWORD_FILE"))
	password := strings.TrimRight(string(passwordBytes), "\r\n")
	database := os.Getenv("POSTGRES_DB")

	switch stl {
	case "mem":
		lst = linkmemstore.NewLinks()
	case "pg":
		dsn := fmt.Sprintf("postgresql://%v:%v@postgres/%v?sslmode=disable", username, password, database)
		pgst, err := pgstore.NewLinks(dsn)
		if err != nil {
			log.Fatal(err)
		}
		defer pgst.Close()
		lst = pgst
	default:
		log.Panic("unknown URL_SHORTENER_STORE = ", stl)
	}

	app := starter.NewApp(lst)
	links := link.NewLinks(lst)

	h := handler.NewRouter(links)

	srv := server.NewServer(":8000", h)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
