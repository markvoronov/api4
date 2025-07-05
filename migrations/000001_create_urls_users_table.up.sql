CREATE TABLE users (
                       id   BIGSERIAL     PRIMARY KEY,
                       name VARCHAR       NOT NULL UNIQUE
);

CREATE TABLE urls (
                      id            BIGSERIAL       PRIMARY KEY,
                      url           VARCHAR         NOT NULL UNIQUE,
                      alias         VARCHAR         NOT NULL UNIQUE,
                      created_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
                      created_user  BIGINT
                          REFERENCES users(id)

);