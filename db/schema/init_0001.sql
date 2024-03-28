CREATE SCHEMA
    IF NOT EXISTS rate_api;

CREATE TABLE
    IF NOT EXISTS rate_api.exchange_rates (
        id VARCHAR(255) PRIMARY KEY,
        date DATE NOT NULL,
        currency CHAR(3) NOT NULL,
        rate DECIMAL(10, 4) NOT NULL
);

CREATE INDEX
    IF NOT EXISTS exchange_rates_date_currency_idx
    ON rate_api.exchange_rates (date, currency);