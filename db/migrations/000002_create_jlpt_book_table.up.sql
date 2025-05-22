CREATE TABLE IF NOT EXISTS jlpt_books (
    jlpt_book_id UUID PRIMARY KEY,
    created_at BIGINT NOT NULL,
    modified_at BIGINT,
    deleted_at BIGINT,
    created_by VARCHAR NOT NULL,
    modified_by VARCHAR,
    deleted_by VARCHAR,
    name VARCHAR NOT NULL,
    level VARCHAR NOT NULL,
    category VARCHAR,
    year VARCHAR,
    source_type VARCHAR NOT NULL,
    url VARCHAR
);