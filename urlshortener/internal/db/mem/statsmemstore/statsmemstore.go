package statsmemstore

import (
	"context"
	"sync"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/stats"
	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

var _ stats.StatsStoreInterface = &Stats{}

type Stats struct {
	sync.Mutex
	m map[shorturl.ShortUrl]stats.Stats
}

func (s *Stats) Create(ctx context.Context, stats stats.Stats) error {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:

	}
	s.m[stats.Short] = stats
	return nil
}

func (s *Stats) GetByLink(ctx context.Context, short shorturl.ShortUrl) (chan stats.Stats, error) {
	s.Lock()
	defer s.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	chout := make(chan stats.Stats, 100)

	go func() {
		defer close(chout)
		s.Lock()
		defer s.Unlock()
		for i, v := range s.m {
			if i == short {
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):

				case chout <- v:

				}
			}
		}
	}()
	return chout, nil
}
