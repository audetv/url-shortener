CREATE TABLE public.links (
	short varchar NOT NULL,
	url text NOT NULL,
	search text,
	redirect_count int,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL,
    CONSTRAINT links_pk PRIMARY KEY (short)
);

CREATE INDEX links_url_idx ON public.links (url);
CREATE INDEX links_search_idx ON public.links (search);
CREATE INDEX links_redirect_count_idx ON public.links (redirect_count);
