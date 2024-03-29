CREATE SCHEMA
    IF NOT EXISTS rate_api;

-- Table: rate_api.exchange_rates
-- Composite key: date, currency
-- Indexes: date, currency, date+currency


CREATE TABLE
    IF NOT EXISTS rate_api.exchange_rates (
        -- id VARCHAR(36) NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
        day DATE NOT NULL,
        -- day_str CHAR(10) NOT NULL,
        currency CHAR(3) NOT NULL,
        rate DECIMAL(10, 4) NOT NULL,
        PRIMARY KEY (day, currency)
);