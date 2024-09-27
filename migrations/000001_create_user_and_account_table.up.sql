CREATE TABLE accounts (
    id VARCHAR PRIMARY KEY,
    username VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    role VARCHAR NOT NULL,
    created_at DATE,
    updated_at DATE
);

CREATE TABLE users (
    id VARCHAR PRIMARY KEY,
    name VARCHAR NOT NULL,
    avatar VARCHAR,
    email VARCHAR,
    account_id VARCHAR,
    created_at DATE,
    updated_at DATE,
    CONSTRAINT account_fk FOREIGN KEY (account_id) REFERENCES accounts(id)
);