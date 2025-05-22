CREATE TABLE IF NOT EXISTS questions (
    question_id UUID PRIMARY KEY,
    created_at BIGINT NOT NULL,
    modified_at BIGINT,
    deleted_at BIGINT,
    created_by VARCHAR NOT NULL,
    modified_by VARCHAR,
    deleted_by VARCHAR,
    quiz_id UUID NOT NULL REFERENCES quizzes(quiz_id) ON DELETE CASCADE,
    question_text VARCHAR NOT NULL,
    question_type VARCHAR
);