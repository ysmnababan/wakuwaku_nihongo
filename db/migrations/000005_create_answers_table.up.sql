CREATE TABLE IF NOT EXISTS answers (
    answer_id UUID PRIMARY KEY,
    created_at BIGINT NOT NULL,
    modified_at BIGINT,
    deleted_at BIGINT,
    created_by VARCHAR NOT NULL,
    modified_by VARCHAR,
    deleted_by VARCHAR,
    question_id UUID NOT NULL REFERENCES questions(question_id) ON DELETE CASCADE,
    answer_text VARCHAR NOT NULL,
    is_correct BOOLEAN NOT NULL 
);