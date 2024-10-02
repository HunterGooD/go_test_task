-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS musics (
    id SERIAL PRIMARY KEY,
    m_group VARCHAR(60) NOT NULL,
    m_song VARCHAR(100) NOT NULL,
    m_link VARCHAR(100) NOT NULL,
    m_text TEXT,
    m_release_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_text  ON musics USING gin (m_text gin_trgm_ops);

CREATE INDEX IF NOT EXISTS idx_group ON musics (m_group);
CREATE INDEX IF NOT EXISTS idx_song  ON musics (m_song);
CREATE INDEX IF NOT EXISTS idx_link  ON musics (m_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS musics CASCADE;
-- +goose StatementEnd
