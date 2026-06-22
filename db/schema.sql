CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    secret_hash BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_activated_at TIMESTAMP,
    total_activations INTEGER NOT NULL,
    blacklisted BOOLEAN NOT NULL
);