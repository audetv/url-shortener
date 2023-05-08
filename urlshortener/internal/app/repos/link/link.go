package link

import (
	"context"
	"fmt"
	"net/url"
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
	Create(ctx context.Context, link Link) (*Link, error)
	Read(ctx context.Context, short shorturl.ShortUrl) (*Link, error)
	SearchLinks(ctx context.Context, short string) (chan Link, error)
	IncRedirectCount(ctx context.Context, short shorturl.ShortUrl) error
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

func (ls *Links) CreateLink(ctx context.Context, l Link) (*Link, error) {
	l.Short = *shorturl.New(8)

	l = decodeLink(l)

	newLink, err := ls.linkStore.Create(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}

	l.Short = newLink.Short
	l.CreatedAt = newLink.CreatedAt
	return &l, nil
}

func (ls *Links) Read(ctx context.Context, short shorturl.ShortUrl) (*Link, error) {
	link, err := ls.linkStore.Read(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("read link error %w", err)
	}
	return link, err
}

func (ls *Links) DoRedirect(ctx context.Context, short shorturl.ShortUrl) (*Link, error) {
	var link, err = ls.linkStore.Read(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ссылки %w", err)
	}

	err = ls.linkStore.IncRedirectCount(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("ошибка при увелечении счетчика переходов по ссылке %w", err)
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
	return chout, nil
}

// decodeLink декодирует оригинальную ссылку, если произошла ошибка возвращает не изменённую ссылку
// это нужно, чтобы был возможен текстовый поиск по ссылке
func decodeLink(link Link) Link {
	decodedValue, err := url.QueryUnescape(link.Origin)
	if err != nil {
		decodedValue = link.Origin
	}

	link.Origin = decodedValue
	return link
}
