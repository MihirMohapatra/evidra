-- +goose Up
CREATE TABLE questions (
    id               UUID         PRIMARY KEY,
    questionnaire_id UUID         NOT NULL REFERENCES questionnaires(id) ON DELETE CASCADE,
    text             TEXT         NOT NULL,
    type             VARCHAR(30)  NOT NULL DEFAULT 'open_ended',
    "order"          INT          NOT NULL DEFAULT 0,
    options          TEXT[]       NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_questionnaire ON questions(questionnaire_id);

-- +goose Down
DROP TABLE IF EXISTS questions;
