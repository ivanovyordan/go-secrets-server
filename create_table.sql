CREATE TABLE secrets (
    hash CHAR(64) UNIQUE PRIMARY KEY NOT NULL,
    secret_text text NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    remaining_views INTEGER NOT NULL
);
