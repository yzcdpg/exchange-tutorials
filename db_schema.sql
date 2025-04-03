-- db_schema.sql
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    mode TEXT DEFAULT 'isolated'
);

CREATE TABLE IF NOT EXISTS trades (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    price TEXT NOT NULL,
    quantity TEXT NOT NULL,
    buyer_id TEXT NOT NULL,
    seller_id TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS balances (
    user_id TEXT REFERENCES users(id),
    currency TEXT NOT NULL,
    balance TEXT NOT NULL,
    PRIMARY KEY (user_id, currency)
);

