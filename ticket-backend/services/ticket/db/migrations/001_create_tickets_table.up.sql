CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('open', 'pending', 'done')),
    notes TEXT
);
