-- Dev init: schema and table for music service (owner: postgres)
CREATE SCHEMA IF NOT EXISTS music;

CREATE TABLE IF NOT EXISTS music.albums (
    id integer GENERATED ALWAYS AS IDENTITY,
    title text NOT NULL,
    artist text NOT NULL,
    price numeric(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT albums_pkey PRIMARY KEY (id)
);
