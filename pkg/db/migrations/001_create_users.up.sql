CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    username VARCHAR(15) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password TEXT NOT NULL,
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_email_key UNIQUE (email)
);

CREATE INDEX users_deleted_at_idx ON users (deleted_at);
