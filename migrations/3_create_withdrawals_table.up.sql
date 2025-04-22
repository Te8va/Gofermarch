BEGIN;

CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    login TEXT NOT NULL REFERENCES users(login),
    order_number TEXT NOT NULL,
    withdrawn DOUBLE PRECISION NOT NULL,
    processed_at TIMESTAMP NOT NULL
);

COMMIT;
