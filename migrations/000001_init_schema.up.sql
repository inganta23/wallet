CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00 CHECK (balance >= 0)
);

INSERT INTO users (username, balance) VALUES ('alice', 1000.00);