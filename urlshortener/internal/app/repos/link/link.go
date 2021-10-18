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
	Log           []Stats
}

type Stats struct {
	Referer   string
	Location  string
	CreatedAt time.Time
}

type LinkStoreInterface interface {
	CreateLink(ctx context.Context, link Link) (*shorturl.ShortUrl, error)
	ReadLink(ctx context.Context, short shorturl.ShortUrl) (*Link, error)
	SearchLinks(ctx context.Context, short string) (chan Link, error)
	IncRedirectCount(ctx context.Context, short shorturl.ShortUrl) (*Link, error)
	// AddStats(ctx context.Context, short shorturl.ShortUrl, stats Stats) error
	// GetStats(ctx context.Context, short shorturl.ShortUrl) (chan Stats, error)
	// Delete(ctx context.Context, su shorturl.ShortUrl) error

}

type Links struct {
	linkStore LinkStoreInterface
}

func NewLinks(linkStore LinkStoreInterface) *Links {
	return &Links{
		linkStore: linkStore,
	}
}

func (ls *Links) DoRedirect(ctx context.Context, short shorturl.ShortUrl, stats Stats) (*Link, error) {
	link, err := ls.linkStore.IncRedirectCount(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("read link error %w", err)
	}
	return link, err
}

func (ls *Links) CreateLink(ctx context.Context, l Link) (*Link, error) {
	l.Short = *shorturl.New()
	l.CreatedAt = time.Now()
	short, err := ls.linkStore.CreateLink(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}
	l.Short = *short
	return &l, nil
}

func (ls *Links) Read(ctx context.Context, short shorturl.ShortUrl) (*Link, error) {
	link, err := ls.linkStore.ReadLink(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("read link error %w", err)
	}
	return link, err
}

func (ls *Links) SearchLinks(ctx context.Context, s string) (chan Link, error) {
	chin, err := ls.linkStore.SearchLinks(ctx, s)
	if err != nil {
		return nil, err
	}

	chout := make(chan Link, 100)
	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case l, ok := <-chin:
				if !ok {
					return
				}
				chout <- l
			}
		}
	}()
	return chout, err
}
