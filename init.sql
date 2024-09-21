CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id VARCHAR(255) NOT NULL PRIMARY KEY,
    token BYTEA NOT NULL,
    ip_address VARCHAR(45) NOT NULL
);
