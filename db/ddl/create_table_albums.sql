-- Table: music.albums

-- DROP TABLE IF EXISTS music.albums;

CREATE TABLE IF NOT EXISTS music.albums
(
    id integer NOT NULL DEFAULT nextval('music.albums_id_seq'::regclass),
    title text COLLATE pg_catalog."default" NOT NULL,
    artist text COLLATE pg_catalog."default" NOT NULL,
    price numeric(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT albums_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS music.albums
    OWNER to ryandayrit;
