CREATE SCHEMA
    IF NOT EXISTS rate_api;

CREATE TABLE
    IF NOT EXISTS rate_api.exchange_rates (
        date DATE PRIMARY KEY,
        rates JSONB
);