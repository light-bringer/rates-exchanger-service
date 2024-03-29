package cron

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// PeriodicJob runs the given task function periodically with the specified interval
// PeriodicJob runs a cron job with the specified task and interval.
// The cron job will be started in a goroutine and can be stopped gracefully
// by sending an OS signal (Interrupt or SIGTERM).
//
// Parameters:
//   - task: A function that represents the task to be performed by the cron job.
//   - interval: The time duration between each execution of the task.
//
// Example usage:
//
//	PeriodicJob(func() {
//	  // Perform the task here
//	}, 1 * time.Hour)
//
// Note: The cron job will continue running until an OS signal is received.
func PeriodicJob(task func(), interval time.Duration) {
	// Create a channel to receive signals for stopping the cron job
	stop := make(chan struct{})

	// Handle OS signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Start the cron job in a goroutine
	go func() {
		defer close(stop)
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-stop:
				slog.Debug("Stopping cron job...")
				ticker.Stop()
				return
			case <-ticker.C:
				task()
			}
		}
	}()

	// Wait for OS signals to gracefully shut down the cron job
	<-signalCh
	slog.Info("Received termination signal. Stopping cron job...")
	close(stop)

	// Wait for the goroutine to stop
	time.Sleep(1 * time.Second)
	slog.Info("Cron job stopped.")
}

func Periodically(ctx context.Context, task func(), interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // Listen for context cancellation to gracefully stop the job
			slog.Debug("Stopping cron job...")
			return
		case <-ticker.C:
			task()
		}
	}
}
