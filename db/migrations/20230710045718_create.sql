-- migrate:up
CREATE TABLE github_account (
    internal_id CHAR(26) PRIMARY KEY,
    loginname VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    FOREIGN KEY FK_GITHUBACCOUNT_INTERNALID (internal_id)
    REFERENCES users(id)
);
UNIQUE INDEX INDEX_EXTERNAL_ACCOUNT_INTERNAL_ID ON github_account;
INDEX INDEX_EXTERNAL_ACCOUNT_LOGINNAME ON github_account;

-- migrate:down

DROP UNIQUE INDEX INDEX_USERS_ID ON github_account;
DROP TABLE github_account;
