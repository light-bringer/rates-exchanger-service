package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // pgx driver
	"github.com/light-bringer/rates-exchanger-service/models"
	"github.com/pkg/errors"
)

const (
	pgxMaxConnLifetime   = 60 * time.Minute
	pgxMaxConnIdleTime   = 15 * time.Minute
	pgxHealthCheckPeriod = 1 * time.Hour
)

// BuildPGXConnPool returns a new pgx connection pool.
func BuildPGXConnPool(ctx context.Context, pgConf models.PostgresConfig) (*pgxpool.Pool, error) {
	sqlAddr := buildPostgresConnString(pgConf)
	connConfig, err := pgxpool.ParseConfig(sqlAddr)
	if err != nil {
		slog.Error("Error parsing connection string")
		return nil, errors.Wrap(err, "error parsing connection string")
	}

	connConfig.MaxConns = int32(pgConf.MaxConnections)
	connConfig.MinConns = int32(pgConf.MinConnections)
	connConfig.MaxConnLifetime = pgxMaxConnLifetime
	connConfig.MaxConnIdleTime = pgxMaxConnIdleTime
	connConfig.HealthCheckPeriod = pgxHealthCheckPeriod

	connPool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		slog.Error("Error creating connection pool", "err", err)
		return nil, errors.Wrap(err, "error creating connection pool")
	}

	if _, err = connPool.Exec(ctx, "SELECT 1"); err != nil {
		slog.Error("Failed to ping database", "err", err)
		return nil, errors.Wrap(err, "failed to ping database")
	}
	slog.Info("Connected to postgres")

	if !checkIfTableExistsInDatabase(connPool) {
		return nil, errors.New("table does not exist in database, please re-run migration again")
	}

	return connPool, nil
}

// checkIfTableExistsInDatabase check if table exists in database
// and return false if it does not.
func checkIfTableExistsInDatabase(db *pgxpool.Pool) bool {
	selectBuilder := squirrel.Select("table_name").
		From("information_schema.tables").
		Where(squirrel.Eq{"table_name": TableName}).
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := selectBuilder.ToSql()
	if err != nil {
		slog.Error("Unable to build SQL query", "error", err)
		return false
	}

	slog.Info("Checking if table exists", "sql", sql, "args", args)

	results, err := db.Query(context.Background(), sql, args...)
	if err != nil {
		slog.Error("query execution error", "error", err)
		return false
	}

	for results.Next() {
		var tableName string
		if err := results.Scan(&tableName); err != nil {
			slog.Error("Failed to scan row", "error", err)
			return false
		}
		if tableName != TableName {
			return false
		}
	}
	return true
}

// buildPostgresConnString returns a postgres connection string.
func buildPostgresConnString(pgConf models.PostgresConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pgConf.Host,
		pgConf.Port,
		pgConf.Username,
		pgConf.Password,
		pgConf.Database,
	)
}
