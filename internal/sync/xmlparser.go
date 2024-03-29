package sync

import (
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/pkg/errors"
)

type ExchangeRateSync struct {
	httpClient *http.Client
	url        string
	schema     string
	tableName  string
	db         *pgxpool.Pool
}

func NewExchangeRateSync(schemaName, url string, db *pgxpool.Pool) *ExchangeRateSync {
	if schemaName == "" {
		schemaName = "public"
	}
	return &ExchangeRateSync{
		httpClient: http.DefaultClient,
		url:        url,
		db:         db,
		schema:     schemaName,
		tableName:  schemaName + ".exchange_rates",
	}
}

// loadHTTPData loads the exchange rates from the given URL.
func (e *ExchangeRateSync) loadHTTPData() (models.ExchangeRates, error) {
	req, reqErr := http.NewRequestWithContext(context.Background(), http.MethodGet, e.url, nil)
	if reqErr != nil {
		slog.Error("Error creating request", reqErr)
		return models.ExchangeRates{}, errors.Wrap(reqErr, "error creating request")
	}

	resp, err := e.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("Error getting exchange rates", err)
		return models.ExchangeRates{}, errors.Wrap(err, "error getting exchange rates")
	}

	defer resp.Body.Close()
	var envelope models.Envelope

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body", err)
		return models.ExchangeRates{}, errors.Wrap(err, "error reading response body")
	}
	if xmlErr := xml.Unmarshal(body, &envelope); xmlErr != nil {
		return models.ExchangeRates{}, errors.Wrap(xmlErr, "error unmarshalling response body")
	}

	slog.Info("Exchange rates loaded successfully", "count", len(envelope.Cube.Cubes))

	res := make(models.ExchangeRates, 0, len(envelope.Cube.Cubes))
	for _, cube := range envelope.Cube.Cubes {
		cubeTime, _ := time.Parse("2006-01-02", cube.Time)
		for _, entry := range cube.Entries {
			rate, _ := strconv.ParseFloat(entry.Rate, 64)
			res = append(res, models.ExchangeRate{
				Currency: entry.Currency,
				Rate:     rate,
				Time:     cubeTime,
			})
		}
	}

	slog.Info("Exchange rates parsed successfully", "count", len(res), "rates", res[:5])
	return res, nil
}

// insertToDB inserts the exchange rates into the database using a transaction and batch inserts.
func (e *ExchangeRateSync) insertToDB(exchangeRates models.ExchangeRates) error {
	ctx := context.Background()
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	tx, err := e.db.Begin(ctx)
	if err != nil {
		slog.Error("Error beginning transaction", err)
		return errors.Wrap(err, "error beginning transaction")
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(context.Background())
			if rollbackErr != nil {
				slog.Error("Error rolling back transaction", "error", rollbackErr)
			}
			return
		}
		commitErr := tx.Commit(context.Background())
		if commitErr != nil {
			slog.Error("Error committing transaction", "error", commitErr)
		}
	}()

	insertQueryBuilder := sq.Insert(e.tableName).Columns("currency", "rate", "day")

	batchSize := 1000
	for idx := 0; idx < len(exchangeRates); idx += batchSize {
		batchEnd := idx + batchSize
		if batchEnd > len(exchangeRates) {
			batchEnd = len(exchangeRates)
		}

		batchRates := exchangeRates[idx:batchEnd]
		if err := e.insertBatchToDB(tx, insertQueryBuilder, batchRates); err != nil {
			return err
		}

		// Reset insertQueryBuilder for the next batch
		insertQueryBuilder = sq.Insert(e.tableName).Columns("currency", "rate", "day")
	}

	slog.Info("Exchange rates inserted successfully")
	return nil
}

// insertBatchToDB inserts a batch of exchange rates into the database.
func (e *ExchangeRateSync) insertBatchToDB(
	tx pgx.Tx,
	insertQueryBuilder squirrel.InsertBuilder,
	batchRates models.ExchangeRates,
) error {
	for _, rate := range batchRates {
		insertQueryBuilder = insertQueryBuilder.Values(rate.Currency, rate.Rate, rate.Time)
	}
	insertQueryBuilder = insertQueryBuilder.Suffix(`ON CONFLICT (day, currency) DO UPDATE SET rate = EXCLUDED.rate`)
	query, args, queryErr := insertQueryBuilder.ToSql()
	if queryErr != nil {
		slog.Error("Error building insert query", queryErr)
		return errors.Wrap(queryErr, "error building insert query")
	}

	res, execErr := tx.Exec(context.Background(), query, args...)
	if execErr != nil {
		slog.Error("Error inserting exchange rates", execErr)
		return errors.Wrap(execErr, "error inserting exchange rates")
	}

	slog.Debug("Inserted exchange rates", "rows", res.RowsAffected(), "query", query, "args", args)
	slog.Info("Inserted exchange rates", "rows", res.RowsAffected())

	return nil
}

// Sync synchronizes the exchange rates with the external API.
func (e *ExchangeRateSync) Sync() {
	exchangeRates, err := e.loadHTTPData()
	slog.Debug("Exchange rates loaded", "exchangeRates", exchangeRates, "error", err)
	if err != nil {
		slog.Error("Error loading exchange rates", err)
		return
	}

	if insertErr := e.insertToDB(exchangeRates); insertErr != nil {
		slog.Info("Exchange rates synchronized successfully")
	}
}

// deleteOldRates deletes the exchange rates older than the specified number of days.
func (e *ExchangeRateSync) deleteOldRates(days int) error {
	ctx := context.Background()
	sq := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	threshold := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	deleteBuilder := sq.Delete(e.tableName).Where(squirrel.Gt{"day": threshold})

	query, args, queryErr := deleteBuilder.ToSql()
	if queryErr != nil {
		slog.Error("Error building delete query", queryErr)
		return errors.Wrap(queryErr, "error building delete query")
	}

	res, execErr := e.db.Exec(ctx, query, args...)
	if execErr != nil {
		slog.Error("Error deleting exchange rates", execErr)
		return errors.Wrap(execErr, "error deleting exchange rates")
	}

	slog.Debug("Deleted old exchange rates", "rows", res.RowsAffected(), "query", query, "args", args)
	slog.Info("Deleted old exchange rates", "rows", res.RowsAffected())

	return nil
}

// Cleanup deletes the exchange rates older than the specified number of days.
func (e *ExchangeRateSync) Cleanup(days int) {
	if err := e.deleteOldRates(days); err != nil {
		slog.Error("Error cleaning up exchange rates", err)
	}
}
