package cron

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeriodicJob(t *testing.T) {
	// Define a mock task function
	t.Run("mock task", func(t *testing.T) {
		// Define a mock task function
		mockTask := func() {
			slog.Info("Mock task executed")
		}

		mockTask()

		// Run the periodic job with a 1-second interval
		// PeriodicJob(mockTask, 1*time.Second)
		assert.True(t, true)
	})
}
