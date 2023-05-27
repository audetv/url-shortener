package pgstore

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // Postgresql driver

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
	"github.com/audetv/url-shortener/urlshortener/internal/app/shorturl"
)

var _ link.LinkStoreInterface = &Links{}

type DBPgLink struct {
	Short         string     `db:"short"`
	Url           string     `db:"url"`
	Search        string     `db:"search"`
	RedirectCount int        `db:"redirect_cont"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

type Links struct {
	db *sql.DB
}

func NewLinks(dsn string) (*Links, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	ls := &Links{
		db: db,
	}
	return ls, nil
}

func (ls *Links) Close() {
	ls.db.Close()
}

func (ls *Links) Create(ctx context.Context, l link.Link) (*link.Link, error) {
	dbl := &DBPgLink{
		Short:         string(l.Short),
		Url:           l.Origin,
		Search:        l.Search,
		RedirectCount: l.RedirectCount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := ls.db.ExecContext(ctx, `INSERT INTO links
    (short, url, search, redirect_count, created_at, updated_at, deleted_at)
    values ($1, $2, $3, $4, $5, $6, $7)`,
		dbl.Short,
		dbl.Url,
		dbl.Search,
		dbl.RedirectCount,
		dbl.CreatedAt,
		dbl.UpdatedAt,
		nil,
	)
	if err != nil {
		return nil, err
	}

	l.CreatedAt = dbl.CreatedAt
	return &l, nil
}

func (ls *Links) Delete(ctx context.Context, short shorturl.ShortUrl) error {
	_, err := ls.db.ExecContext(ctx, `UPDATE links SET deleted_at = $2 WHERE short = $1`,
		short, time.Now(),
	)
	return err
}

func (ls *Links) Read(ctx context.Context, short shorturl.ShortUrl) (*link.Link, error) {
	dbl := &DBPgLink{}
	rows, err := ls.db.QueryContext(ctx, `SELECT short, url, search, redirect_count, created_at, updated_at, deleted_at
	FROM links WHERE short = $1`, short)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&dbl.Short,
			&dbl.Url,
			&dbl.Search,
			&dbl.RedirectCount,
			&dbl.CreatedAt,
			&dbl.UpdatedAt,
			&dbl.DeletedAt,
		); err != nil {
			return nil, err
		}
	}

	return &link.Link{
		Short:         *shorturl.Parse(dbl.Short),
		Origin:        dbl.Url,
		Search:        dbl.Search,
		RedirectCount: dbl.RedirectCount,
		CreatedAt:     dbl.CreatedAt,
	}, nil
}

// ReadByOrigin ищет ссылку в БД по совпадению url
func (ls *Links) ReadByOrigin(ctx context.Context, l link.Link) (*link.Link, error) {
	dbl := &DBPgLink{}
	rows, err := ls.db.QueryContext(ctx, `SELECT short, url, search, redirect_count, created_at, updated_at, deleted_at
	FROM links WHERE url = $1`, l.Origin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&dbl.Short,
			&dbl.Url,
			&dbl.Search,
			&dbl.RedirectCount,
			&dbl.CreatedAt,
			&dbl.UpdatedAt,
			&dbl.DeletedAt,
		); err != nil {
			return nil, err
		}
	}

	return &link.Link{
		Short:         *shorturl.Parse(dbl.Short),
		Origin:        dbl.Url,
		Search:        dbl.Search,
		RedirectCount: dbl.RedirectCount,
		CreatedAt:     dbl.CreatedAt,
	}, nil
}

func (ls *Links) IncRedirectCount(ctx context.Context, short shorturl.ShortUrl) error {
	_, err := ls.db.ExecContext(
		ctx,
		`UPDATE links SET redirect_count = redirect_count + 1, updated_at = $1 WHERE short = $2`,
		time.Now(),
		short,
	)
	return err
}

func (ls *Links) SearchLinks(ctx context.Context, s string) (chan link.Link, error) {
	chout := make(chan link.Link, 100)

	go func() {
		defer close(chout)
		dbl := &DBPgLink{}

		rows, err := ls.db.QueryContext(ctx, `
		SELECT short, url, search, redirect_count, created_at, updated_at, deleted_at 
		FROM links WHERE search LIKE $1  ORDER BY redirect_count DESC`, "%"+s+"%")
		if err != nil {
			log.Println(err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(
				&dbl.Short,
				&dbl.Url,
				&dbl.Search,
				&dbl.RedirectCount,
				&dbl.CreatedAt,
				&dbl.UpdatedAt,
				&dbl.DeletedAt,
			); err != nil {
				log.Println(err)
				return
			}

			chout <- link.Link{
				Short:         *shorturl.Parse(dbl.Short),
				Origin:        dbl.Url,
				Search:        dbl.Search,
				RedirectCount: dbl.RedirectCount,
				CreatedAt:     dbl.CreatedAt,
			}
		}
	}()

	return chout, nil
}
