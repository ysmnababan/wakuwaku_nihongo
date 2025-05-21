CREATE TABLE IF NOT EXISTS quizzes (
    quiz_id UUID PRIMARY KEY,
    created_at BIGINT NOT NULL,
    modified_at BIGINT,
    deleted_at BIGINT,
    created_by VARCHAR NOT NULL,
    modified_by VARCHAR,
    deleted_by VARCHAR,
    title VARCHAR NOT NULL,
    description VARCHAR
);