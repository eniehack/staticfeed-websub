-- migrate:up
CREATE TABLE sites (
    id CHAR(26) PRIMARY KEY,
    host URL NOT NULL,
    feed_url URL NOT NULL,
    owned_by CHAR(26) NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY FK_SITES_OWNED_BY (owned_by)
    REFERENCES users(id)
);
CREATE INDEX INDEX_SITES_HOST ON sites (host);

-- migrate:down
DROP INDEX INDEX_SITES_HOST ON sites;
DROP TABLE sites;