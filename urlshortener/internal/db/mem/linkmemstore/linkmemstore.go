package linkmemstore

import (
	"context"
	"sync"

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

func (ls *Links) Create(ctx context.Context, l link.Link) (*shorturl.ShortUrl, error) {
	ls.Lock()
	defer ls.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	ls.m[l.Short] = l
	return &l.Short, nil
}
