-- +goose Up
-- +goose StatementBegin
CREATE TABLE watchers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_email TEXT NOT NULL,
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_watchers_route ON watchers(origin, destination);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS watchers;
-- +goose StatementEnd
