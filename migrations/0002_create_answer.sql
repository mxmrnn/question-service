-- +goose Up
CREATE TABLE IF NOT EXISTS answers (
    id           BIGSERIAL PRIMARY KEY,
    question_id  BIGINT      NOT NULL,
    user_id      VARCHAR(64) NOT NULL,
    text         TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_answers_question
    FOREIGN KEY (question_id)
    REFERENCES questions (id)
    ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_answers_question_id ON answers (question_id);

-- +goose Down
DROP TABLE IF EXISTS answers;