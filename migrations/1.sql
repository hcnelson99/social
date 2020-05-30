BEGIN;

DROP TABLE IF EXISTS version;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE version(
    lock char(1) NOT NULL DEFAULT 'X',
    current int NOT NULL,
    constraint pk_t1 PRIMARY KEY (lock),
    constraint ck_t1_locked CHECK (lock='X')
);

INSERT INTO version(current) VALUES(1);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    text TEXT NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    session_generation INTEGER DEFAULT 0,
    username TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL
    -- disabling/banning users?
);

COMMIT;
