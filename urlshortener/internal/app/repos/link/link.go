package link

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

type Link struct {
	Short         shorturl.ShortUrl
	Search        string
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
	ReadByOrigin(ctx context.Context, link Link) (*Link, error)
	SearchLinks(ctx context.Context, s string) (chan Link, error)
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

	l = processLink(l)

	// Еси ссылка уже существует, найдена по origin, то возвращаем найденную ссылку
	existLink, err := ls.checkExistsLink(ctx, l)
	if err == nil {
		return existLink, nil
	}

	newLink, err := ls.linkStore.Create(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("create link error: %w", err)
	}

	l.Short = newLink.Short
	l.CreatedAt = newLink.CreatedAt
	return &l, nil
}

func (ls *Links) checkExistsLink(ctx context.Context, l Link) (*Link, error) {
	result, err := ls.linkStore.ReadByOrigin(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("read link error %w", err)
	}
	if result.Short == "" {
		return nil, fmt.Errorf("not found")
	}
	return result, nil
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
	// приводим поисковый запрос в нижней регистр перед запросом в БД
	s = strings.ToLower(s)
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

// processLink подготавливает ссылку перед сохранением в БД
// декодирует и парсит UrlQuery, разбирает часть запроса "search[query]" и приводит к нижнему регистру
func processLink(l Link) Link {
	// Декодирование ссылки надо делать перед тем как применить strings.ToLower
	l = decodeLink(l)
	u, err := url.Parse(l.Origin)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()

	var result string
	for n, search := range q["search[query]"] {
		if n == 0 {
			result = fmt.Sprintf("%v", strings.TrimSpace(search))
		} else {
			result = fmt.Sprintf("%v, %v", result, strings.TrimSpace(search))
		}
	}

	fmt.Println(result)
	l.Search = strings.ToLower(result)

	return l
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
