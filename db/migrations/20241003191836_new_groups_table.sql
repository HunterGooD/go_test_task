-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS groups (
    id SERIAL PRIMARY KEY,
    g_name VARCHAR(60) NOT NULL,
    g_link VARCHAR(100) NOT NULL,
    g_start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_group ON groups (g_name);
CREATE INDEX IF NOT EXISTS idx_group_link  ON groups (g_link);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
-- +goose StatementEnd
