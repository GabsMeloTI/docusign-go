CREATE TABLE contracts (
    id SERIAL PRIMARY KEY,
    provider_name VARCHAR(255) NOT NULL,
    document_url TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    is_signed BOOL,
    date_signed TIMESTAMP
);


