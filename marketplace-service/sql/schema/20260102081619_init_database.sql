-- +goose Up
-- +goose StatementBegin
CREATE TABLE trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    driver_id UUID NOT NULL,
    price_thb INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'OPEN',
    departure_time TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trips_origin_dest ON trips(origin, destination);
CREATE INDEX idx_trips_status ON trips(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trips;
-- +goose StatementEnd
