-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    g_name VARCHAR(60) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_group ON groups (g_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
