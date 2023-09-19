-- +goose Up

-- theme
CREATE TABLE theme (
	id      INTEGER PRIMARY KEY,
	title   TEXT NOT NULL,
	img_url TEXT
);

-- theme_question
CREATE TABLE theme_question (
	id       INTEGER PRIMARY KEY,
	theme_id INTEGER NOT NULL,
	music_id INTEGER NOT NULL,
	text     TEXT NOT NULL,
	hint     TEXT
);

-- music
CREATE TABLE music (
	id        INTEGER PRIMARY KEY,
	deezer_id INTEGER NOT NULL,
	artist_id INTEGER,
	album_id  INTEGER,
	name      TEXT NOT NULL,
	mp3_url   TEXT NOT NULL
);

-- music_artist
CREATE TABLE music_artist (
	id        INTEGER PRIMARY KEY,
	deezer_id INTEGER NOT NULL,
	name      TEXT NOT NULL,
	img_url   TEXT
);

-- music_album
CREATE TABLE music_album (
	id        INTEGER PRIMARY KEY,
	deezer_id INTEGER NOT NULL,
	name      TEXT NOT NULL,
	img_url   TEXT
);

-- +goose Down

DROP TABLE music_artist;
DROP TABLE music_album;
DROP TABLE music;
DROP TABLE theme_question;
DROP TABLE theme;