CREATE TABLE public.links (
	short varchar NOT NULL,
	url varchar(512) NOT NULL,
	redirect_count int,
	created_at timestamptz NOT NULL,
	updated_at timestamptz NOT NULL,
	deleted_at timestamptz NULL
)
