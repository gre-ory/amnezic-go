-- +goose Up

-- user
CREATE TABLE user (
	id      	INTEGER PRIMARY KEY,
	name    	TEXT NOT NULL,
	hash 		TEXT NOT NULL,
	permissions TEXT
);

-- user session
CREATE TABLE user_session (
	token		TEXT PRIMARY KEY,
	user_id    	INTEGER,
	expiration	INTEGER
);

-- +goose Down

DROP TABLE user_session;
DROP TABLE user;