CREATE TABLE IF NOT EXISTS customers (
    customer_id UUID PRIMARY KEY,
    created_at BIGINT NOT NULL,
    modified_at BIGINT,
    deleted_at BIGINT,
    created_by VARCHAR NOT NULL,
    modified_by VARCHAR,
    deleted_by VARCHAR,

    username VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR,
    is_active BOOLEAN NOT NULL
);