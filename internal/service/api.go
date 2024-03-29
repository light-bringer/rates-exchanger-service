package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/pkg/errors"
)

// FetchLatestExchangeRates fetches the latest exchange rates.
// The function returns the exchange rates for the latest day.
// The rates are sorted in ascending order.
// The function returns an error if the query fails.
func (s *RatesService) FetchLatestExchangeRates(limit uint64) (models.LatestExchangeRates, error) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	subQuery := queryBuilder.Select("MAX(day) AS latest_day").From(s.tableName)

	subQueryStr, _, _ := subQuery.ToSql()

	// Main query: Joins the exchange_rates table with the subquery to fetch rates for the latest day.
	query := queryBuilder.Select("er.currency", "er.rate").
		From(s.tableName + " AS er").
		JoinClause(fmt.Sprintf("INNER JOIN (%s) AS ld ON er.day = ld.latest_day", subQueryStr)).
		OrderBy("er.rate ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		slog.Error("Failed to build SQL query", "error", err)
		return nil, errors.Wrap(err, "failed to build SQL query")
	}

	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		slog.Error("Failed to acquire connection", "error", err)
		return nil, errors.Wrap(err, "failed to acquire connection")
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), sqlStr, args...)
	if err != nil {
		slog.Error("Failed to execute query", "error", err)
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	rates := make(models.LatestExchangeRates, 0)

	// Iterate through the result set.
	for rows.Next() {
		var currency string
		var rate float64
		if err := rows.Scan(&currency, &rate); err != nil {
			slog.Error("Failed to scan row", "error", err)
			continue
		}

		rates = append(rates, models.LatestExchangeRate{
			Currency: currency,
			Rate:     rate,
		})
	}

	return rates, nil
}

// FetchRatesForDate fetches the exchange rates for a given date.
// The function returns the exchange rates for the given date.
// The rates are sorted in ascending order.
// The function returns an error if the query fails.
func (s *RatesService) FetchRatesForDate(
	date string,
	limit uint64,
) (models.LatestExchangeRates, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		slog.Error("Failed to parse date", "error", err)
		return nil, errors.Wrap(err, "failed to parse date")
	}

	selectBuilder := psql.Select("currency", "rate").From(s.tableName).Where(squirrel.Eq{"day": date}).OrderBy("currency ASC")

	if limit > 0 {
		selectBuilder = selectBuilder.Limit(limit)
	}

	sqlStr, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL query")
	}

	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		slog.Error("Failed to acquire connection", "error", err)
		return nil, errors.Wrap(err, "failed to acquire connection")
	}

	defer conn.Release()

	rows, err := conn.Query(context.Background(), sqlStr, args...)
	if err != nil {
		slog.Error("Failed to execute query", "error", err)
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	rates := make(models.LatestExchangeRates, 0)

	// Iterate through the result set.
	for rows.Next() {
		var currency string
		var rate float64
		if err := rows.Scan(&currency, &rate); err != nil {
			slog.Error("Failed to scan row", "error", err)
			continue
		}

		rates = append(rates, models.LatestExchangeRate{
			Currency: currency,
			Rate:     rate,
		})
	}

	return rates, nil
}

// GetRateStatistics fetches the rate statistics for the latest day.
// The statistics include min, max, and average rates for each currency.
// The statistics are calculated based on the rates for the latest day.
// The function returns a map of currency to rate statistics.
func (s *RatesService) GetRateStatistics(days uint64) (models.RateStatisticsMap, error) {
	/**
	SELECT
	    currency,
	    MIN(rate) AS min_rate,
	    MAX(rate) AS max_rate,
	    AVG(rate) AS avg_rate
	FROM
	    rate_api.exchange_rates
	WHERE
	    day = (SELECT MAX(day) FROM rate_api.exchange_rates)
	GROUP BY
	    currency
	ORDER BY
	    currency;
		**/
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	subQuery := psql.Select("MAX(day)").From(s.tableName)
	subQueryStr, _, _ := subQuery.ToSql()

	whereClause := fmt.Sprintf("day <= (%s) AND day >= (%s) - INTERVAL '%d days'", subQueryStr, subQueryStr, days)

	query := psql.Select(
		"currency",
		"MIN(rate) AS min_rate",
		"MAX(rate) AS max_rate",
		"AVG(rate) AS avg_rate",
	).From("rate_api.exchange_rates").
		Where(whereClause).
		GroupBy("currency").
		OrderBy("currency")

	sqlStr, args, err := query.ToSql()
	slog.Info("SQL query", "sql", sqlStr, "args", args)
	if err != nil {
		slog.Error("Failed to build SQL query", "error", err)
		return nil, errors.Wrap(err, "failed to build SQL query")
	}

	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		slog.Error("Failed to acquire connection", "error", err)
		return nil, errors.Wrap(err, "failed to acquire connection")
	}
	defer conn.Release()

	rows, err := conn.Query(context.Background(), sqlStr, args...)
	if err != nil {
		slog.Error("Failed to execute query", "error", err)
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	stats := make(models.RateStatisticsMap)
	for rows.Next() {
		var stat models.RateStatistic
		if err := rows.Scan(&stat.Currency, &stat.MinRate, &stat.MaxRate, &stat.AvgRate); err != nil {
			slog.Error("Failed to read row", "error", err)
			continue
		}
		stats[stat.Currency] = stat
	}

	return stats, nil
}
