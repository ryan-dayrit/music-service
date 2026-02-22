-- Integration test init script for PostgreSQL
-- Creates schema and albums table without user-specific clauses

CREATE SCHEMA IF NOT EXISTS music;

CREATE TABLE IF NOT EXISTS music.albums
(
    id integer GENERATED ALWAYS AS IDENTITY,
    title text NOT NULL,
    artist text NOT NULL,
    price numeric(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT albums_pkey PRIMARY KEY (id)
);
