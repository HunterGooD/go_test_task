-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_update_at_column() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN NEW.update_at = (now() at time zone 'UTC'); RETURN NEW; END; $$;

CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    m_name VARCHAR(100) NOT NULL,
    m_link VARCHAR(100) NOT NULL,
    m_text TEXT,
    m_release_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE OR REPLACE TRIGGER update_data_systems_modtime BEFORE UPDATE ON songs FOR EACH ROW EXECUTE PROCEDURE update_update_at_column();

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_text  ON songs USING gin (m_text gin_trgm_ops) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_song  ON songs (m_name) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_link  ON songs (m_link) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_song_deleted_at  ON songs (coalesce(deleted_at));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
DROP FUNCTION IF EXISTS update_update_at_column;
-- +goose StatementEnd
