package service

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockDBConn simulates a database connection for testing.
type MockDBConn struct {
	mock.Mock
}

func (m *MockDBConn) Query(ctx context.Context, sql string, args ...interface{}) (*Rows, error) {
	m.Called(ctx, sql, args)
	// Simplified: Return a fictional Rows object and nil error. Adapt as needed.
	return &Rows{}, nil
}

// Rows simulates database query results.
type Rows struct {
	mock.Mock
	currentIndex int
	data         []map[string]interface{} // Simulate your data structure
}

func (r *Rows) Next() bool {
	r.currentIndex++
	return r.currentIndex < len(r.data)
}

func (r *Rows) Close() error {
	return nil // Simplify: Assume no error
}

func (r *Rows) Scan(dest ...interface{}) error {
	// Simplified scan, you'd need to populate dest based on your data
	return nil
}
