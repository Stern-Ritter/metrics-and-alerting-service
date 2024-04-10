DROP TABLE metrics;
DROP TABLE metric_types;

CREATE TABLE metric_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL
);

CREATE UNIQUE INDEX metric_type_name_idx ON metric_types (name);

CREATE TABLE metrics (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    type_id INTEGER NOT NULL,
    value DOUBLE PRECISION NOT NULL,
    FOREIGN KEY (type_id) REFERENCES metric_types (id)
);

CREATE INDEX metric_name_idx ON metrics (name);

INSERT INTO metric_types (name) VALUES ('gauge'), ('counter');