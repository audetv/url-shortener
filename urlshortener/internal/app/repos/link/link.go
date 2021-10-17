package link

import (
	"context"
	"fmt"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

type Link struct {
	Short         shorturl.ShortUrl
	Origin        string
	RedirectCount int
	CreatedAt     time.Time
}

type Stats struct {
	Short     shorturl.ShortUrl
	Referrer  string
	Location  string
	CreatedAt time.Time
}

type LinkStoreInterface interface {
	Create(ctx context.Context, l Link) (*shorturl.ShortUrl, error)
	// Read(ctx context.Context, su shorturl.ShortUrl) (*Link, error)
	// Delete(ctx context.Context, su shorturl.ShortUrl) error
	// SearchLinks(ctx context.Context, su string) (chan Link, error)
	// GetStats(ctx context.Context, su shorturl.ShortUrl) (*Stats, error)
	// CreateStats(ctx context.Context, su shorturl.ShortUrl) (*Stats, error)
	// IncRedirectCount(ctx context.Context, su shorturl.ShortUrl) (*Stats, error)
}

type Links struct {
	linkStore LinkStoreInterface
}

func NewLinks(linkStore LinkStoreInterface) *Links {
	return &Links{
		linkStore: linkStore,
	}
}

func (ls *Links) CreateLink(ctx context.Context, l Link) (*Link, error) {
	l.Short = *shorturl.New()
	short, err := ls.linkStore.Create(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}
	l.Short = *short
	return &l, nil
}
