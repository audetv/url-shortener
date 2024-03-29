package linkmemstore

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

var _ link.LinkStoreInterface = &Links{}

type Links struct {
	sync.Mutex
	m map[shorturl.ShortUrl]link.Link
}

func NewLinks() *Links {
	return &Links{
		m: make(map[shorturl.ShortUrl]link.Link),
	}
}

func (ls *Links) Create(ctx context.Context, l link.Link) (*link.Link, error) {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ls.m[l.Short] = l
	return &l, nil
}

func (ls *Links) Read(ctx context.Context, short shorturl.ShortUrl) (*link.Link, error) {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}
	l, ok := ls.m[short]
	if ok {
		return &l, nil
	}

	return nil, sql.ErrNoRows
}

// ReadByOrigin эта функция не работает, TODO реализовать поиск по origin
func (ls *Links) ReadByOrigin(ctx context.Context, l link.Link) (*link.Link, error) {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}
	l, ok := ls.m[l.Short]
	if ok {
		return &l, nil
	}

	return nil, sql.ErrNoRows
}

func (ls *Links) IncRedirectCount(ctx context.Context, short shorturl.ShortUrl) error {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}

	l, ok := ls.m[short]
	if ok {
		l.RedirectCount += 1
		ls.m[short] = l
		return nil
	}

	return sql.ErrNoRows
}

func (ls *Links) SearchLinks(ctx context.Context, s string) (chan link.Link, error) {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	chout := make(chan link.Link, 100)

	go func() {
		defer close(chout)
		ls.Lock()
		defer ls.Unlock()
		for _, l := range ls.m {
			if strings.Contains(l.Origin, s) {
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):

				case chout <- l:
				}
			}
		}
	}()
	return chout, nil
}
