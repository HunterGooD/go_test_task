-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    g_name VARCHAR(60) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE OR REPLACE TRIGGER update_data_systems_modtime BEFORE UPDATE ON groups FOR EACH ROW EXECUTE PROCEDURE update_update_at_column();

CREATE INDEX IF NOT EXISTS idx_group ON groups (g_name)  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_group_deleted_at  ON groups (coalesce(deleted_at));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
