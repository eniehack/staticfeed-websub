-- migrate:up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    expired_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX INDEX_USERS_ID ON users (id);

-- migrate:down
DROP UNIQUE INDEX INDEX_USERS_ID ON users;
DROP TABLE users;

