package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

type Stats struct {
	Short     shorturl.ShortUrl
	Referrer  string
	Location  string
	CreatedAt time.Time
}

type StatsStoreInterface interface {
	Create(ctx context.Context, stats Stats) error
	GetByLink(ctx context.Context, short shorturl.ShortUrl) (chan Stats, error)
}

type StatsLog struct {
	statsStore StatsStoreInterface
}

func NewStats(statsStore StatsStoreInterface) *StatsLog {
	return &StatsLog{
		statsStore: statsStore,
	}
}

func (sl *StatsLog) Create(ctx context.Context, stats Stats) error {
	err := sl.statsStore.Create(ctx, stats)
	if err != nil {
		return fmt.Errorf("create stats log error: %w", err)
	}
	return nil
}

func (sl *StatsLog) GetByLink(ctx context.Context, short shorturl.ShortUrl) (chan Stats, error) {
	chin, err := sl.statsStore.GetByLink(ctx, short)
	if err != nil {
		return nil, err
	}
	chout := make(chan Stats, 100)

	go func() {
		defer close(chout)
		for {
			select {
			case <-ctx.Done():
				return
			case s, ok := <-chin:
				if !ok {
					return
				}
				chout <- s
			}
		}
	}()
	return chout, nil
}
