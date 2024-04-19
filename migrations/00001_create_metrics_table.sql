-- +goose Up
-- +goose StatementBegin
CREATE TABLE metric_types
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

CREATE TABLE metrics
(
    id    BIGSERIAL PRIMARY KEY,
    name  VARCHAR(256)     NOT NULL,
    type  METRIC_TYPE      NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

CREATE INDEX metric_name_idx ON metrics (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
DROP TYPE metric_type;
-- +goose StatementEnd