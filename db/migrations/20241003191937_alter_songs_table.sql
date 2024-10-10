-- +goose Up
-- +goose StatementBegin
ALTER TABLE songs
  ADD COLUMN IF NOT EXISTS group_id INTEGER REFERENCES groups(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE songs
  DROP COLUMN IF EXISTS group_id;
-- +goose StatementEnd
