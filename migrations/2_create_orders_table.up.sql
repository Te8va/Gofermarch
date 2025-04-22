BEGIN;

CREATE TABLE IF NOT EXISTS orders (
    number TEXT PRIMARY KEY,
    login TEXT NOT NULL REFERENCES users(login),
    uploaded_at TIMESTAMP NOT NULL,
    status TEXT DEFAULT 'NEW' CHECK (status IN ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED')),
    accrual DOUBLE PRECISION
);

COMMIT;
